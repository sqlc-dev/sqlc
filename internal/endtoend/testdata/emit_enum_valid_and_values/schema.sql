CREATE TYPE ip_protocol AS enum ('tcp', 'ip', 'icmp');

CREATE TABLE bar_old (id_old serial not null, ip_old ip_protocol NOT NULL);

