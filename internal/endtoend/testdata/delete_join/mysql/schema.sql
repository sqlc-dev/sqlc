CREATE TABLE primary_table (
        id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        user_id bigint(20) unsigned NOT NULL,
        PRIMARY KEY (id)
);

CREATE TABLE join_table (
        id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        primary_table_id bigint(20) unsigned NOT NULL,
        other_table_id bigint(20) unsigned NOT NULL,
        is_active tinyint(1) NOT NULL DEFAULT '0',
        PRIMARY KEY (id)
);

