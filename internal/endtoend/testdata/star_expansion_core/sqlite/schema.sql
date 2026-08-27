CREATE TABLE foo (a text, b text, "group" text);
CREATE TABLE bar (a text, c text);
CREATE TABLE baz (z text);
CREATE TABLE recipes (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL);
CREATE VIRTUAL TABLE recipes_fts USING FTS5 (name, content='recipes', content_rowid='id');
