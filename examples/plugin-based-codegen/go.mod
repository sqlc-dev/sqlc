module github.com/sqlc-dev/sqlc/examples/plugin-based-codegen

go 1.24.0

require (
	github.com/sqlc-dev/sqlc v1.30.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/grpc v1.77.0 // indirect
)

// Use local sqlc for development
replace github.com/sqlc-dev/sqlc => ../..
