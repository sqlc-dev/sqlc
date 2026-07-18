CREATE TABLE parent_table (
    id bigserial PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE child_table (
    id bigserial PRIMARY KEY,
    parent_id bigint NOT NULL REFERENCES parent_table (id),
    x integer NOT NULL,
    y text NOT NULL
);
