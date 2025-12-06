package golang

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
)

func TestClickHouseDriver_NotPGX(t *testing.T) {
	driver := opts.SQLDriver(opts.SQLDriverClickHouseV2)
	if driver.IsPGX() {
		t.Error("ClickHouse driver should not identify as PGX")
	}
}

func TestClickHouseDriver_IsClickHouse(t *testing.T) {
	driver := opts.SQLDriver(opts.SQLDriverClickHouseV2)
	if !driver.IsClickHouse() {
		t.Error("ClickHouse driver should identify as ClickHouse")
	}
}

func TestStandardDriver_NotClickHouse(t *testing.T) {
	driver := opts.SQLDriver(opts.SQLDriverLibPQ)
	if driver.IsClickHouse() {
		t.Error("Standard driver should not identify as ClickHouse")
	}
}
