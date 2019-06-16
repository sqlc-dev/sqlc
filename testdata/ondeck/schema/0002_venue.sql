CREATE TABLE venue (
    slug             text primary key,
    name             text not null,
    city             text references city(slug),
    spotify_playlist text not null,
    songkick_id      text
)
