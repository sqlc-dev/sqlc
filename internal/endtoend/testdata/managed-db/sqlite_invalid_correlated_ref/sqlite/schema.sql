CREATE TABLE organizations (
    id INTEGER PRIMARY KEY
);

CREATE TABLE organization_members (
    id INTEGER PRIMARY KEY,
    organization_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL
);

CREATE TABLE projects (
    id INTEGER PRIMARY KEY,
    organization_id INTEGER NOT NULL
);

CREATE TABLE locations (
    id INTEGER PRIMARY KEY,
    public_id TEXT UNIQUE NOT NULL,
    project_id INTEGER NOT NULL
) STRICT;
