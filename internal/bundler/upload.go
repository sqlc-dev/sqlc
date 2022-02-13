package bundler

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/kyleconroy/sqlc/internal/config"
)

type Uploader struct {
	configPath string
	config     *config.Config
}

func NewUploader(configPath string, conf *config.Config) *Uploader {
	return &Uploader{
		configPath: configPath,
		config:     conf,
	}
}

func (up *Uploader) Upload(ctx context.Context) error {
	body := bytes.NewBuffer([]byte{})

	gw := gzip.NewWriter(body)
	defer gw.Close()
	w := multipart.NewWriter(gw)
	defer w.Close()

	if err := writeContents(w, up.configPath, up.config); err != nil {
		return err
	}

	if err := gw.Flush(); err != nil {
		return err
	}
	w.Close()
	gw.Close()

	req, err := http.NewRequest("POST", "http://localhost:8000/upload", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("code=%d", resp.StatusCode)
	}
	return nil
}
