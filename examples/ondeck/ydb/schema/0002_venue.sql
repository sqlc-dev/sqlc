CREATE TABLE venues (
    id Serial NOT NULL,
    dropped Text,
    status Text NOT NULL,
    slug Text NOT NULL,
    name Text NOT NULL,
    city Text NOT NULL,
    spotify_playlist Text NOT NULL,
    songkick_id Text,
    tags Text,
    PRIMARY KEY (id)
);

