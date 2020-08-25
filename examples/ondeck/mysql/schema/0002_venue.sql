CREATE TABLE venues (
    id               SERIAL primary key,
    dropped          text,
    status           ENUM('open', 'closed') not null COMMENT 'Venues can be either open or closed',
    statuses         text, -- status[],
    slug             text         not null COMMENT 'This value appears in public URLs',
    name             varchar(255) not null,
    city             text         not null references city(slug),
    spotify_playlist varchar(255) not null,
    songkick_id      text,
    tags             text  -- text[]
) COMMENT='Venues are places where muisc happens';
