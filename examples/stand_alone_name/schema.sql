create table ding_depts
(
    id    bigint not null
        constraint ding_depts_pk
            primary key,
    pid   bigint,
    title varchar
);

comment on table ding_depts is '钉钉部门';

comment on column ding_depts.id is '部门id';

comment on column ding_depts.pid is '上级部门id';

comment on column ding_depts.title is '部门名称';

create index ding_depts_pid_index
    on ding_depts (pid);

create table domains
(
    tag     varchar                                               not null
        constraint domain_pk
            primary key,
    leaders character varying[] default '{}'::character varying[] not null,
    configs jsonb               default '{}'::jsonb               not null
);

comment on table domains is '领域';

comment on column domains.tag is '领域标签';

comment on column domains.leaders is '领域领导';

comment on column domains.configs is '领域配置';

create index domain_leaders_index
    on domains (leaders);
