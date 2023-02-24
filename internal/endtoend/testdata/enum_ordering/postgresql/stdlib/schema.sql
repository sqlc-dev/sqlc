CREATE TYPE enum_type AS ENUM ('first', 'last');
ALTER TYPE enum_type ADD VALUE 'third' AFTER 'first';
ALTER TYPE enum_type ADD VALUE 'fourth' BEFORE 'last';
ALTER TYPE enum_type ADD VALUE 'fifth' AFTER 'fourth';
ALTER TYPE enum_type ADD VALUE 'second' BEFORE 'third';

CREATE TABLE foo (
    id SERIAL PRIMARY KEY
);