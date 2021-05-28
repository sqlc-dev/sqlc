CREATE SCHEMA foo;
CREATE TABLE foo.bar (baz text);
CREATE TYPE foo.bat AS ENUM ('bat');
COMMENT ON SCHEMA foo IS 'Schema comment';
COMMENT ON TABLE foo.bar IS 'Table comment';
COMMENT ON COLUMN foo.bar.baz IS 'Column comment';
COMMENT ON TYPE foo.bat IS 'Enum comment';