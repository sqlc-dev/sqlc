CREATE TABLE venues (
    id               SERIAL primary key,
    dropped          text,
    slug             text not null,
    name             text not null,
    city             text references city(slug),
    spotify_playlist text not null,
    songkick_id      text
)
