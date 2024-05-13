-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA content;
CREATE TABLE "content"."order" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "number" varchar NOT NULL,
    "user_id" uuid NOT NULL,
    "created_at" timestamptz NOT NULL,
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."order";
DROP SCHEMA content;
-- +goose StatementEnd
