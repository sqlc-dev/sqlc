package dinosql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const inputMigration = `
-- +goose Up
ALTER TABLE archived_jobs ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE;

-- +goose Down
ALTER TABLE archived_jobs DROP COLUMN expires_at;
`

const outputMigration = `
-- +goose Up
ALTER TABLE archived_jobs ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE;
`

func TestRemoveGooseRollback(t *testing.T) {
	if diff := cmp.Diff(outputMigration, RemoveGooseRollback(inputMigration)); diff != "" {
		t.Errorf("migration mismatch:\n%s", diff)
	}
}
