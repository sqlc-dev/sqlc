-- Test ALTER TABLE operations with mixed case table and column names
-- Verifies consistent handling of case sensitivity in DDL operations
CREATE TABLE Users (id integer primary key, name text, "Email" text);

-- Test renaming columns with different case formats
ALTER TABLE Users RENAME COLUMN name TO full_name;
ALTER TABLE Users RENAME COLUMN "Email" TO "EmailAddress";

-- Test adding a simple column
ALTER TABLE Users ADD COLUMN created_at text;
