package bundler

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

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

func (up *Uploader) Upload(ctx context.Context, result map[string]string) error {
	body := bytes.NewBuffer([]byte{})

	// gw := gzip.NewWriter(body)
	// defer gw.Close()
	w := multipart.NewWriter(body)
	defer w.Close()

	if err := writeInputs(w, up.configPath, up.config); err != nil {
		return err
	}
	if err := writeOutputs(w, up.dir, result); err != nil {
		return err
	}

	// if err := gw.Flush(); err != nil {
	// 	return err
	// }
	w.Close()
	// gw.Close()

	req, err := http.NewRequest("POST", "http://localhost:8090/upload", body)
	if err != nil {
		return err
	}

	// Set sqlc-version header
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
