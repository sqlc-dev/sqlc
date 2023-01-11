-- Network Address Types
-- https://www.postgresql.org/docs/current/datatype-net-types.html
CREATE TABLE dt_net_types (
  a  inet,
  b  cidr,
  c  macaddr
);

CREATE TABLE dt_net_types_not_null (
  a  inet NOT NULL,
  b  cidr NOT NULL,
  c  macaddr NOT NULL
);
