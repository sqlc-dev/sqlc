CREATE TABLE city (
    slug text PRIMARY KEY,
    name text NOT NULL
);

CREATE TYPE status AS ENUM ('open', 'closed');

CREATE TABLE venue (
    id               SERIAL primary key,
    create_at        timestamp    not null,
    status           status       not null,
    slug             text         not null,
    name             varchar(255) not null,
    city             text         not null references city(slug),
    spotify_playlist varchar      not null,
    songkick_id      text
);
