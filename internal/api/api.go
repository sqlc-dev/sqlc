// Package api is intended to be the future public API for sqlc.
//
// The shape of this package is inspired by esbuild's Build API
// (https://pkg.go.dev/github.com/evanw/esbuild/pkg/api#hdr-Build_API): a small
// surface area of options structs and result structs that lets callers drive
// sqlc programmatically without going through the CLI.
//
// Today the package lives under internal/ while the API stabilises. Once the
// surface settles it is expected to graduate to pkg/api so it can be imported
// by external Go programs.
package api
