package bundler

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/kyleconroy/sqlc/internal/config"
)

type Uploader struct {
	configPath string
	config     *config.Config
	dir        string
}

func NewUploader(configPath, dir string, conf *config.Config) *Uploader {
	return &Uploader{
		configPath: configPath,
		config:     conf,
		dir:        dir,
	}
}

func (up *Uploader) Validate() error {
	if up.config.Project.ID == "" {
		return fmt.Errorf("project ID is not set")
	}
	return nil
}

func (up *Uploader) Upload(ctx context.Context, result map[string]string) error {
	if err := up.Validate(); err != nil {
		return err
	}
	body := bytes.NewBuffer([]byte{})

	w := multipart.NewWriter(body)
	defer w.Close()

	if err := writeInputs(w, up.configPath, up.config); err != nil {
		return err
	}
	if err := writeOutputs(w, up.dir, result); err != nil {
		return err
	}

	w.Close()

	req, err := http.NewRequest("POST", "http://localhost:8090/upload", body)
	if err != nil {
		return err
	}

	// Set sqlc-version header
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SQLC_AUTH_TOKEN")))
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("upload endpiont returned non-200 status code: %d", resp.StatusCode)
	}
	return nil
}
