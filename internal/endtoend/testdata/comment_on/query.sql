CREATE SCHEMA foo;

CREATE TABLE foo.bar (
        baz text NOT NULL 
);

CREATE TYPE foo.mood AS ENUM ('sad', 'ok', 'happy');

COMMENT ON SCHEMA foo IS 'this is the foo schema';
COMMENT ON TYPE foo.mood IS 'this is the mood type';
COMMENT ON TABLE foo.bar IS 'this is the bar table';
COMMENT ON COLUMN foo.bar.baz IS 'this is the baz column';

-- name: ListBar :many
SELECT * FROM foo.bar;
