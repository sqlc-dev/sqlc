create table descriptions (
    id varchar(32) GENERATED ALWAYS AS (MD5(txt)) STORED,
    txt text,
    primary key (id)
);
