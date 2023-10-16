create table attribute_value
(
    id              bigserial not null,
    val             text      not null,
    attribute       bigint    not null
);

create table attribute
(
    id              bigserial      not null,
    name            text           not null
);
