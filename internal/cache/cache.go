package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

// The cache directory defaults to os.UserCacheDir(). This location can be
// overridden by the SQLCCACHE environment variable.
//
// Currently the cache stores two types of data: plugins and query analysis
func Dir() (string, error) {
	cache := os.Getenv("SQLCCACHE")
	if cache != "" {
		return cache, nil
	}
	cacheHome, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheHome, "sqlc"), nil
}

func PluginsDir() (string, error) {
	cacheRoot, err := Dir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cacheRoot, "plugins")
	if err := os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("failed to create %s directory: %w", dir, err)
	}
	return dir, nil
}

func AnalysisDir() (string, error) {
	cacheRoot, err := Dir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cacheRoot, "query_analysis")
	if err := os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("failed to create %s directory: %w", dir, err)
	}
	return dir, nil
}
