create table RouterGroup
(
    groupId                int unsigned auto_increment primary key,
    groupName              varchar(100)         not null,
    defaultConfigId        int unsigned         null,
    defaultFirmwareVersion varchar(12)          null,
    parentGroupId int unsigned      null,
    firmwarePolicy         varchar(45)          null,
    styles                 text                 null
);
