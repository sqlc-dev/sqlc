-- +goose Up
-- +goose StatementBegin
CREATE TABLE "auth"."user" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "login" varchar NOT NULL,
    "password" varchar NOT NULL,
    "created_at" timestamptz NOT NULL,
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "auth"."user";
-- +goose StatementEnd
