CREATE TABLE venues (
    id               integer primary key AUTOINCREMENT,
    dropped          text,
    status           text not null,
    statuses         text, -- status[]
    slug             text         not null,
    name             varchar(255) not null,
    city             text         not null references city(slug),
    spotify_playlist varchar(255) not null,
    songkick_id      text,
    tags             text, -- tags[]
    CHECK (status = 'open' OR status = 'closed')
);
