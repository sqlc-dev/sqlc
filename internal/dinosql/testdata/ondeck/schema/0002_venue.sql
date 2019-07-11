CREATE TYPE status AS ENUM ('open', 'closed');

CREATE TABLE venues (
    id               SERIAL primary key,
    dropped          text,
    status           status       not null,
    slug             text         not null,
    name             varchar(255) not null,
    city             text         not null references city(slug),
    spotify_playlist varchar      not null,
    songkick_id      text
)
