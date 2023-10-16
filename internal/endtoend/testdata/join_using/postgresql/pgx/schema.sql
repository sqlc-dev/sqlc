create table t1 (
        fk integer not null unique
);
create table t2 (
        fk integer not null references t1(fk)
);
