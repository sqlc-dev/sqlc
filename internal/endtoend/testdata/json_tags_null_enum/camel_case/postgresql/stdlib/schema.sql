CREATE TYPE job_post_location_type AS ENUM('remote', 'in_office', 'hybrid');

CREATE TABLE authors (
    id   BIGSERIAL PRIMARY KEY,
    type job_post_location_type,
    name text      NOT NULL,
    bio  text
);

