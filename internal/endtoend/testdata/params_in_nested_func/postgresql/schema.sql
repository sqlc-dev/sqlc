create table RouterGroup
(
    groupId                serial primary key,
    groupName              varchar(100)         not null,
    defaultConfigId        int         null,
    defaultFirmwareVersion varchar(12)          null,
    parentGroupId int      null,
    firmwarePolicy         varchar(45)          null,
    styles                 text                 null
);
