-- name: InsertRecords :copyfrom
INSERT INTO records (id, kind, optional_kind, note) VALUES (?, ?, ?, ?);

-- name: ReplaceRecords :copyfrom
REPLACE INTO records (id, kind, optional_kind, note) VALUES (?, ?, ?, ?);

-- name: InsertSingleValue :copyfrom
INSERT INTO single_values (kind) VALUES (?);

-- name: InsertNullableValue :copyfrom
INSERT INTO nullable_values (kind) VALUES (?);

-- name: InsertOverriddenValue :copyfrom
INSERT INTO overridden_values (id, kind) VALUES (?, ?);

-- name: InsertTextValue :copyfrom
INSERT INTO text_values (value) VALUES (?);
