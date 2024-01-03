CREATE TABLE `Calendar` (
  `Id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `Relation` bigint(20) unsigned NOT NULL,
  `CalendarName` longblob NOT NULL,
  `Title` longblob NOT NULL,
  `Description` longblob NOT NULL,
  `Timezone` varchar(50) NOT NULL,
  `UniqueKey` varchar(50) NOT NULL,
  `IdKey` varchar(50) NOT NULL,
  `MainCalendar` enum('true','false') NOT NULL DEFAULT 'false',
  `CreateDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ModifyDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`Id`),
  KEY `Relation` (`Relation`),
  KEY `UniqueKey` (`UniqueKey`),
  KEY `IdKey` (`IdKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `Event` (
  `Id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `Relation` bigint(20) unsigned NOT NULL,
  `CalendarReference` bigint(20) unsigned NOT NULL,
  `UniqueKey` varchar(50) NOT NULL,
  `EventName` longblob NOT NULL,
  `Description` longblob NOT NULL,
  `Location` varchar(500) NOT NULL,
  `Timezone` varchar(50) NOT NULL,
  `IdKey` varchar(48) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  KEY `Relation` (`Relation`),
  KEY `CalendarReference` (`CalendarReference`),
  KEY `UniqueKey` (`UniqueKey`),
  KEY `IdKey` (`IdKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE authors (
  id   BIGINT AUTO_INCREMENT NOT NULL,
  address VARCHAR(200) NOT NULL,
  name VARCHAR(20) NOT NULL,
  bio  LONGTEXT NOT NULL
);
