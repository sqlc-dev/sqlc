CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.accounts (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    state character varying
);

CREATE TABLE public.users_accounts (
    ID2 uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying
);

