ifeq ($(GOBIN),)
GOBIN = $(firstword $(subst :, ,$(GOPATH)))/bin
endif

SQLC := $(GOBIN)/sqlc
SQL_FILES = $(shell find ./query -type f -name '*.sql')

$(SQLC):
	go get github.com/kyleconroy/sqlc/cmd/sqlc

prepared/db.go:
prepared/models.go:
db.go:
models.go: $(SQL_FILES) sqlc.json | $(SQLC)
	$(SQLC) generate

sql: db.go models.go prepared/models.go prepared/db.go
