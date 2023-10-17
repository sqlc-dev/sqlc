CREATE TABLE "user" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "username" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "telephone" int NOT NULL DEFAULT 0,
  "default_payment" bigint,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "address" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "address_line" varchar NOT NULL,
  "region" varchar NOT NULL,
  "city" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "user_address" (
  "user_id" bigint NOT NULL,
  "address_id" bigint UNIQUE NOT NULL,
  "default_address" bigint,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);
