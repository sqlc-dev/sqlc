CREATE TABLE "mytable" (
    id                BIGSERIAL   NOT NULL PRIMARY KEY,
    myjson             JSONB       NOT NULL
);

insert into mytable (myjson) values
    ('{}'),
    ('{"thing1": {"thing2": "thing3"}}');