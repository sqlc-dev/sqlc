CREATE TYPE enum_type AS ENUM ('first', 'last');
ALTER TYPE enum_type ADD VALUE 'afterlast' AFTER 'last';
ALTER TYPE enum_type ADD VALUE 'third' AFTER 'first';
ALTER TYPE enum_type ADD VALUE 'fourth' BEFORE 'last';
ALTER TYPE enum_type ADD VALUE 'fifth' AFTER 'fourth';
ALTER TYPE enum_type ADD VALUE 'second' BEFORE 'third';
ALTER TYPE enum_type ADD VALUE 'beforefirst' BEFORE 'first';

CREATE TABLE foo (
    id SERIAL PRIMARY KEY
);