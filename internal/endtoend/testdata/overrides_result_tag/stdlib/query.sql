CREATE TABLE public.accounts (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    state character varying
);

CREATE TABLE public.users_accounts (
    ID2 uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying
);

-- name: FindAccount :one
SELECT
    a.*,
    ua.name
    -- other fields
FROM
    accounts a
    INNER JOIN users_accounts ua ON a.id = ua.id2
WHERE
    a.id = @account_id;