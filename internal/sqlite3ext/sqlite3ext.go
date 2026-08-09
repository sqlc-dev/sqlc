// Package sqlite3ext registers the SQLite extensions sqlc relies on with
// github.com/ncruces/go-sqlite3.
//
// As of go-sqlite3 v0.35.0 the FTS5 and R*Tree/Geopoly extensions are compiled
// separately and are no longer part of the default build. Registering them as
// auto extensions makes them available on every connection, matching the
// behavior of earlier releases.
//
// Import this package for side effects wherever a SQLite connection is opened.
package sqlite3ext

import (
	"github.com/ncruces/go-sqlite3"
	"github.com/ncruces/go-sqlite3/ext/fts5"
	"github.com/ncruces/go-sqlite3/ext/rtree"
)

func init() {
	sqlite3.AutoExtension(fts5.Register)
	sqlite3.AutoExtension(rtree.Register)
}
