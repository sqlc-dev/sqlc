CREATE TABLE tbl
(
    value integer
);

-- https://www.postgresql.org/docs/current/sql-createprocedure.html
CREATE PROCEDURE insert_data(
    IN a integer,
    IN b integer,
    -- Numbers
    OUT c integer,
    OUT i float,
    OUT j numeric,
    OUT k real,
    -- Text
    OUT d varchar,
    OUT h text,
    -- Time
    OUT e timestamp,
    OUT m interval,
    -- Other
    OUT f jsonb,
    OUT g bytea,
    OUT l boolean
)
    LANGUAGE plpgsql
AS
$$
BEGIN
    INSERT INTO tbl VALUES (a);
    INSERT INTO tbl VALUES (b);

    c := 777;

    -- Numbers assignments
    i := random() * 100;
    j := (random() * 500)::numeric(10, 2);
    k := (random() * 50)::real;

    -- Text assignments
    d := 'Varchar val ' || floor(random() * 100);
    h := 'Text val ' || md5(random()::text);

    -- Time assignments
    e := now();
    m := make_interval(hours => floor(random() * 24)::int);

    -- Other assignments
    f := '{
      "key": "random"
    }'::jsonb;
    g := '\xDEADBEEF'::bytea;
    l := (random() > 0.5);
END;
$$;
