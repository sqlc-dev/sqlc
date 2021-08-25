CREATE TABLE users (
    id integer,
    first_name varchar(255) NOT NULL
);

CREATE FUNCTION users_func() RETURNS SETOF users AS $func$ BEGIN QUERY
SELECT *
FROM users
END $func$ LANGUAGE plpgsql;
