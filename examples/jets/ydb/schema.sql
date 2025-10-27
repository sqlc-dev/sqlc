CREATE TABLE pilots (
    id Int32 NOT NULL,
    name Text NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE jets (
    id Int32 NOT NULL,
    pilot_id Int32 NOT NULL,
    age Int32 NOT NULL,
    name Text NOT NULL,
    color Text NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE languages (
    id Int32 NOT NULL,
    language Text NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE pilot_languages (
    pilot_id Int32 NOT NULL,
    language_id Int32 NOT NULL,
    PRIMARY KEY (pilot_id, language_id)
);




