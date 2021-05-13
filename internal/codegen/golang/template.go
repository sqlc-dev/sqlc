package golang

import "embed"

//go:embed templates/*
//go:embed templates/pgx/*
//go:embed templates/stdlib/*
var templates embed.FS
