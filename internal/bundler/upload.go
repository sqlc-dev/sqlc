package bundler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/kyleconroy/sqlc/internal/config"
)

type Uploader struct {
	token      string
	configPath string
	config     *config.Config
	dir        string
}

func NewUploader(configPath, dir string, conf *config.Config) *Uploader {
	return &Uploader{
		token:      os.Getenv("SQLC_AUTH_TOKEN"),
		configPath: configPath,
		config:     conf,
		dir:        dir,
	}
}

func (up *Uploader) Validate() error {
	if up.config.Project.ID == "" {
		return fmt.Errorf("project.id is not set")
	}
	if up.token == "" {
		return fmt.Errorf("SQLC_AUTH_TOKEN environment variable is not set")
	}
	return nil
}

func (up *Uploader) buildRequest(ctx context.Context, result map[string]string) (*http.Request, error) {
	body := bytes.NewBuffer([]byte{})
	w := multipart.NewWriter(body)
	if err := writeInputs(w, up.configPath, up.config); err != nil {
		return nil, err
	}
	if err := writeOutputs(w, up.dir, result); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://api.sqlc.dev/upload", body)
	if err != nil {
		return nil, err
	}
	// Set sqlc-version header
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", up.token))
	return req.WithContext(ctx), nil
}

func (up *Uploader) DumpRequestOut(ctx context.Context, result map[string]string) error {
	req, err := up.buildRequest(ctx, result)
	if err != nil {
		return err
	}
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	os.Stdout.Write(dump)
	return nil
}

func (up *Uploader) Upload(ctx context.Context, result map[string]string) error {
	if err := up.Validate(); err != nil {
		return err
	}
	req, err := up.buildRequest(ctx, result)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return fmt.Errorf("upload error: endpoint returned non-200 status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("upload error: %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
