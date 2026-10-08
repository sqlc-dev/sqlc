package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/credentials"
	"oras.land/oras-go/v2/registry/remote/retry"

	"github.com/sqlc-dev/sqlc/internal/info"
)

const (
	wasmLayerMediaType  = "application/wasm"
	wasmConfigMediaType = "application/vnd.wasm.config.v0+json"
)

type wasmConfig struct {
	Architecture string   `json:"architecture"`
	OS           string   `json:"os"`
	LayerDigests []string `json:"layerDigests"`
}

func fetchOCI(ctx context.Context, uri string, plainHTTP bool) ([]byte, error) {
	reference := strings.TrimPrefix(uri, "oci://")
	repo, err := remote.NewRepository(reference)
	if err != nil {
		return nil, fmt.Errorf("parse OCI reference %q: %w", uri, err)
	}
	repo.PlainHTTP = plainHTTP

	store, err := credentials.NewStoreFromDocker(credentials.StoreOptions{})
	if err != nil {
		return nil, fmt.Errorf("load Docker credentials: %w", err)
	}
	client := &auth.Client{
		Client:     retry.DefaultClient,
		Cache:      auth.NewCache(),
		Credential: credentials.Credential(store),
	}
	client.SetUserAgent(fmt.Sprintf("sqlc/%s Go/%s (%s %s)", info.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH))
	repo.Client = client

	desc, manifestReader, err := repo.FetchReference(ctx, repo.Reference.Reference)
	if err != nil {
		return nil, fmt.Errorf("fetch OCI manifest %q: %w", uri, err)
	}
	defer manifestReader.Close()
	if desc.MediaType != ocispec.MediaTypeImageManifest {
		return nil, fmt.Errorf("OCI reference %q resolved to unsupported media type %q", uri, desc.MediaType)
	}
	manifestData, err := content.ReadAll(manifestReader, desc)
	if err != nil {
		return nil, fmt.Errorf("read OCI manifest %q: %w", uri, err)
	}
	var manifest ocispec.Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return nil, fmt.Errorf("decode OCI manifest %q: %w", uri, err)
	}
	if manifest.SchemaVersion != 2 {
		return nil, fmt.Errorf("OCI artifact %q: unsupported schema version %d", uri, manifest.SchemaVersion)
	}
	layer, err := wasmLayer(manifest.Layers)
	if err != nil {
		return nil, fmt.Errorf("OCI artifact %q: %w", uri, err)
	}
	if manifest.Config.MediaType != wasmConfigMediaType {
		return nil, fmt.Errorf("OCI artifact %q: unsupported config media type %q", uri, manifest.Config.MediaType)
	}
	configData, err := content.FetchAll(ctx, repo, manifest.Config)
	if err != nil {
		return nil, fmt.Errorf("fetch OCI WASM config %q: %w", uri, err)
	}
	if err := validateWASMConfig(configData, layer); err != nil {
		return nil, fmt.Errorf("OCI artifact %q: %w", uri, err)
	}

	body, err := repo.Fetch(ctx, layer)
	if err != nil {
		return nil, fmt.Errorf("fetch OCI WASM layer %q: %w", uri, err)
	}
	defer body.Close()
	wmod, err := content.ReadAll(body, layer)
	if err != nil {
		return nil, fmt.Errorf("read OCI WASM layer %q: %w", uri, err)
	}
	return wmod, nil
}

func wasmLayer(layers []ocispec.Descriptor) (ocispec.Descriptor, error) {
	if len(layers) != 1 {
		return ocispec.Descriptor{}, fmt.Errorf("manifest must contain exactly one WASM layer")
	}
	if layers[0].MediaType != wasmLayerMediaType {
		return ocispec.Descriptor{}, fmt.Errorf("unsupported WASM layer media type %q", layers[0].MediaType)
	}
	return layers[0], nil
}

func validateWASMConfig(data []byte, layer ocispec.Descriptor) error {
	var config wasmConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("decode WASM config: %w", err)
	}
	if config.Architecture != "wasm" {
		return fmt.Errorf("unsupported WASM architecture %q", config.Architecture)
	}
	if config.OS != "wasip1" {
		return fmt.Errorf("unsupported WASM OS %q: only wasip1 plugins are supported", config.OS)
	}
	if len(config.LayerDigests) != 1 || config.LayerDigests[0] != layer.Digest.String() {
		return fmt.Errorf("WASM config layerDigests must match the manifest layers")
	}
	return nil
}
