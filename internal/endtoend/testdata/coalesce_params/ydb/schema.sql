CREATE TABLE Calendar (
  Id BigSerial,
  Relation BigSerial,
  CalendarName Text NOT NULL,
  Title Text NOT NULL,
  Description Text NOT NULL,
  Timezone Text NOT NULL,
  UniqueKey Text NOT NULL,
  IdKey Text NOT NULL,
  MainCalendar Bool NOT NULL,
  CreateDate DateTime NOT NULL,
  ModifyDate DateTime NOT NULL,
  PRIMARY KEY (Id)
);

CREATE TABLE Event (
  Id BigSerial,
  Relation BigSerial,
  CalendarReference BigSerial,
  UniqueKey Text NOT NULL,
  EventName Text NOT NULL,
  Description Text NOT NULL,
  Location Text NOT NULL,
  Timezone Text NOT NULL,
  IdKey Text,
  PRIMARY KEY (Id)
);

CREATE TABLE authors (
  id BigSerial,
  address Text NOT NULL,
  name Text NOT NULL,
  bio Text NOT NULL,
  PRIMARY KEY (id)
);
