CREATE TYPE unused_enum AS ENUM (
    'a', 'b'
);

CREATE TYPE unused_table_enum AS ENUM (
    'c', 'd'
);
CREATE TABLE unused_table (
    id      INTEGER PRIMARY KEY,
    value   unused_table_enum
);

CREATE TYPE query_return_full_table_enum AS ENUM (
    'e', 'f'
);
CREATE TABLE query_return_full_table (
    id      INTEGER PRIMARY KEY,
    value   query_return_full_table_enum
);

CREATE TYPE query_param_enum_table_enum AS ENUM (
    'g', 'h'
);
CREATE TABLE query_param_enum_table (
    id      INTEGER PRIMARY KEY,
    other   query_param_enum_table_enum NOT NULL,
    value   query_param_enum_table_enum
);

CREATE TYPE query_param_struct_enum_table_enum AS ENUM (
    'i', 'j'
);
CREATE TABLE query_param_struct_enum_table (
    id      INTEGER PRIMARY KEY,
    value   query_param_struct_enum_table_enum
);

CREATE TYPE query_return_enum_table_enum AS ENUM (
    'k', 'l'
);
CREATE TABLE query_return_enum_table (
    id      INTEGER PRIMARY KEY,
    value   query_return_enum_table_enum
);

CREATE TYPE query_return_struct_enum_table_enum AS ENUM (
    'k', 'l'
);
CREATE TABLE query_return_struct_enum_table (
    id      INTEGER PRIMARY KEY,
    value   query_return_struct_enum_table_enum,
    another INTEGER
);

CREATE TYPE query_sqlc_embed_enum AS ENUM (
    'm', 'n'
);
CREATE TABLE query_sqlc_embed_table (
    id      INTEGER PRIMARY KEY,
    value   query_sqlc_embed_enum
)
