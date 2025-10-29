CREATE TABLE workspaces (
    id uuid NOT NULL,
    owner_id uuid NOT NULL,
    name text NOT NULL
);

CREATE TABLE tasks (
    id uuid NOT NULL,
    workspace_id uuid,
    owner_id uuid NOT NULL,
    name text NOT NULL
);
