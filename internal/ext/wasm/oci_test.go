package wasm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func TestFetchOCIWithDockerCredentials(t *testing.T) {
	wasm := []byte("\x00asm\x01\x00\x00\x00")
	layer := ocispec.Descriptor{
		MediaType: wasmLayerMediaType,
		Digest:    digest.FromBytes(wasm),
		Size:      int64(len(wasm)),
	}
	configData := []byte(fmt.Sprintf(`{"architecture":"wasm","os":"wasip1","layerDigests":[%q]}`, layer.Digest.String()))
	configDigest := digest.FromBytes(configData)
	manifest, err := json.Marshal(ocispec.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		MediaType: ocispec.MediaTypeImageManifest,
		Config: ocispec.Descriptor{
			MediaType: wasmConfigMediaType,
			Digest:    configDigest,
			Size:      int64(len(configData)),
		},
		Layers: []ocispec.Descriptor{layer},
	})
	if err != nil {
		t.Fatal(err)
	}
	manifestDigest := digest.FromBytes(manifest)
	servedManifest := manifest
	servedWASM := wasm
	servedConfig := configData

	const username = "sqlc"
	const password = "secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPassword, ok := r.BasicAuth()
		if !ok || gotUser != username || gotPassword != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="registry"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/v2/plugins/sqlc/manifests/v1", "/v2/plugins/sqlc/manifests/" + manifestDigest.String():
			w.Header().Set("Content-Type", ocispec.MediaTypeImageManifest)
			w.Header().Set("Docker-Content-Digest", manifestDigest.String())
			_, _ = w.Write(servedManifest)
		case "/v2/plugins/sqlc/blobs/" + layer.Digest.String():
			w.Header().Set("Content-Type", wasmLayerMediaType)
			w.Header().Set("Docker-Content-Digest", layer.Digest.String())
			_, _ = w.Write(servedWASM)
		case "/v2/plugins/sqlc/blobs/" + configDigest.String():
			w.Header().Set("Content-Type", wasmConfigMediaType)
			w.Header().Set("Docker-Content-Digest", configDigest.String())
			_, _ = w.Write(servedConfig)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	dockerConfig := t.TempDir()
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	config := fmt.Sprintf(`{"auths":{%q:{"auth":%q}}}`, serverURL.Host, auth)
	if err := os.WriteFile(filepath.Join(dockerConfig, "config.json"), []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DOCKER_CONFIG", dockerConfig)

	// Keep the advertised digests and lengths unchanged so only verification
	// of the response bytes can detect the corruption.
	corruptWASM := bytes.Clone(wasm)
	corruptWASM[len(corruptWASM)-1] ^= 1
	tests := []struct {
		name      string
		reference string
		manifest  []byte
		wasm      []byte
		config    []byte
		wantErr   string
	}{
		{
			name:     "valid artifact",
			manifest: manifest,
			wasm:     wasm,
		},
		{
			name:      "digest reference",
			reference: "@" + manifestDigest.String(),
			manifest:  manifest,
			wasm:      wasm,
		},
		{
			name:     "corrupt manifest",
			manifest: bytes.Replace(manifest, []byte("wasm"), []byte("xasm"), 1),
			wasm:     wasm,
			wantErr:  "read OCI manifest",
		},
		{
			name:     "corrupt config",
			manifest: manifest,
			wasm:     wasm,
			config:   bytes.Replace(configData, []byte("wasip1"), []byte("wasip2"), 1),
			wantErr:  "fetch OCI WASM config",
		},
		{
			name:     "corrupt WASM layer",
			manifest: manifest,
			wasm:     corruptWASM,
			wantErr:  "read OCI WASM layer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			servedManifest = tt.manifest
			servedWASM = tt.wasm
			servedConfig = configData
			if tt.config != nil {
				servedConfig = tt.config
			}
			reference := tt.reference
			if reference == "" {
				reference = ":v1"
			}
			got, err := fetchOCI(context.Background(), "oci://"+serverURL.Host+"/plugins/sqlc"+reference, true)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if gotSum, wantSum := sha256.Sum256(got), sha256.Sum256(wasm); gotSum != wantSum {
				t.Fatalf("WASM checksum = %x, want %x", gotSum, wantSum)
			}
		})
	}
}

func TestWASMLayer(t *testing.T) {
	wasm := ocispec.Descriptor{MediaType: wasmLayerMediaType, Digest: digest.FromString("wasm")}
	other := ocispec.Descriptor{MediaType: "application/octet-stream", Digest: digest.FromString("other")}
	tests := []struct {
		name    string
		layers  []ocispec.Descriptor
		want    ocispec.Descriptor
		wantErr bool
	}{
		{
			name:   "single WASM layer",
			layers: []ocispec.Descriptor{wasm},
			want:   wasm,
		},
		{
			name:    "unsupported layer media type",
			layers:  []ocispec.Descriptor{other},
			wantErr: true,
		},
		{
			name:    "multiple WASM layers",
			layers:  []ocispec.Descriptor{wasm, wasm},
			wantErr: true,
		},
		{
			name:    "WASM layer in multi-layer artifact",
			layers:  []ocispec.Descriptor{other, wasm},
			wantErr: true,
		},
		{
			name:    "no layers",
			wantErr: true,
		},
		{
			name:    "multi-layer artifact without WASM",
			layers:  []ocispec.Descriptor{other, other},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := wasmLayer(tt.layers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wasmLayer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.MediaType != tt.want.MediaType || got.Digest != tt.want.Digest {
				t.Fatalf("wasmLayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateWASMConfig(t *testing.T) {
	layer := ocispec.Descriptor{MediaType: wasmLayerMediaType, Digest: digest.FromString("wasm")}
	tests := []struct {
		name    string
		data    string
		wantErr string
	}{
		{
			name: "WASI Preview 1",
			data: fmt.Sprintf(`{"architecture":"wasm","os":"wasip1","layerDigests":[%q]}`, layer.Digest.String()),
		},
		{
			name:    "invalid JSON",
			data:    `{`,
			wantErr: "decode WASM config",
		},
		{
			name:    "unsupported architecture",
			data:    `{"architecture":"amd64","os":"wasip1"}`,
			wantErr: "unsupported WASM architecture",
		},
		{
			name:    "missing architecture",
			data:    `{"os":"wasip1"}`,
			wantErr: "unsupported WASM architecture",
		},
		{
			name:    "WASI Preview 2 component",
			data:    `{"architecture":"wasm","os":"wasip2","component":{"imports":[],"exports":[]}}`,
			wantErr: "only wasip1 plugins are supported",
		},
		{
			name:    "missing OS",
			data:    `{"architecture":"wasm"}`,
			wantErr: "unsupported WASM OS",
		},
		{
			name:    "missing layer digests",
			data:    `{"architecture":"wasm","os":"wasip1"}`,
			wantErr: "layerDigests must match",
		},
		{
			name:    "mismatched layer digest",
			data:    fmt.Sprintf(`{"architecture":"wasm","os":"wasip1","layerDigests":[%q]}`, digest.FromString("other")),
			wantErr: "layerDigests must match",
		},
		{
			name:    "extra layer digest",
			data:    fmt.Sprintf(`{"architecture":"wasm","os":"wasip1","layerDigests":[%q,%q]}`, layer.Digest, layer.Digest),
			wantErr: "layerDigests must match",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWASMConfig([]byte(tt.data), layer)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
