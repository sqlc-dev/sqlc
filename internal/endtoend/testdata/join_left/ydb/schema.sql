CREATE TABLE users (
    user_id Int32 NOT NULL,
    city_id Int32,
    PRIMARY KEY (user_id)
);

CREATE TABLE cities (
    city_id Int32 NOT NULL,
    mayor_id Int32 NOT NULL,
    PRIMARY KEY (city_id)
);

CREATE TABLE mayors (
    mayor_id Int32 NOT NULL,
    full_name Utf8 NOT NULL,
    PRIMARY KEY (mayor_id)
);

CREATE TABLE authors (
    id Int32 NOT NULL,
    name Utf8 NOT NULL,
    parent_id Int32,
    PRIMARY KEY (id)
);

CREATE TABLE super_authors (
    super_id Int32 NOT NULL,
    super_name Utf8 NOT NULL,
    super_parent_id Int32,
    PRIMARY KEY (super_id)
);

CREATE TABLE users_2 (
    user_id Utf8 NOT NULL,
    user_nickname Utf8 NOT NULL,
    user_email Utf8 NOT NULL,
    user_display_name Utf8 NOT NULL,
    user_password Utf8,
    user_google_id Utf8,
    user_apple_id Utf8,
    user_bio Utf8 NOT NULL DEFAULT '',
    user_created_at Timestamp NOT NULL,
    user_avatar_id Utf8,
    PRIMARY KEY (user_id)
);

CREATE TABLE media (
    media_id Utf8 NOT NULL,
    media_created_at Timestamp NOT NULL,
    media_hash Utf8 NOT NULL,
    media_directory Utf8 NOT NULL,
    media_author_id Utf8 NOT NULL,
    media_width Int32 NOT NULL,
    media_height Int32 NOT NULL,
    PRIMARY KEY (media_id)
);




