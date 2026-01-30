CREATE TABLE routergroup (
    groupId                Serial,
    groupName              Text NOT NULL,
    defaultConfigId        Int32,
    defaultFirmwareVersion Text,
    parentGroupId          Int32,
    firmwarePolicy         Text,
    styles                 Text,
    PRIMARY KEY (groupId)
);

