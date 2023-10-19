package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

// The cache directory defaults to ~/.cache/sqlc. This location can be overriden
// by the SQLCCACHE or XDG_CACHE_HOME environment variable.
//
// Currently the cache stores two types of data: plugins and query analysis
func Dir() (string, error) {
	cache := os.Getenv("SQLCCACHE")
	if cache != "" {
		return cache, nil
	}
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cacheHome = filepath.Join(home, ".cache")
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
