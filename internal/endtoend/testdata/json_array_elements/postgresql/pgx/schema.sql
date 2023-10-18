CREATE TABLE "sys_actions" (
    "id" int8 NOT NULL,
    "code" varchar(20) NOT NULL,
    "menu" varchar(255) NOT NULL,
    "title" varchar(20) NOT NULL,
    "resources" jsonb,
    PRIMARY KEY ("id")
);
