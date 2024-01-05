create table RouterGroup
(
    groupId                int unsigned auto_increment primary key,
    groupName              varchar(100)         not null,
    defaultConfigId        int unsigned         null,
    defaultFirmwareVersion varchar(12)          null,
    inheritPermissions     tinyint(1) default 1 not null,
    parentGroupId int unsigned      null,
    firmwarePolicy         varchar(45)          null,
    styles                 text                 null,
    constraint RouterGroup_ibfk_1
        foreign key (defaultConfigId) references ConfigScript (configId)
            on delete set null
);
