CREATE TABLE venue (
    id               SERIAL primary key,
    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    slug             text not null,
    name             text not null,
    city             text references city(slug),
    spotify_playlist text not null,
    songkick_id      text
)
