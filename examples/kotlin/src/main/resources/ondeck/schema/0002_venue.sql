CREATE TYPE status AS ENUM ('op!en', 'clo@sed');
COMMENT ON TYPE status IS 'Venues can be either open or closed';

CREATE TABLE venues (
    id               SERIAL primary key,
    dropped          text,
    status           status       not null,
    statuses         status[],
    slug             text         not null,
    name             varchar(255) not null,
    city             text         not null references city(slug),
    spotify_playlist varchar      not null,
    songkick_id      text,
    tags             text[]
);
COMMENT ON TABLE venues IS 'Venues are places where muisc happens';
COMMENT ON COLUMN venues.slug IS 'This value appears in public URLs';

