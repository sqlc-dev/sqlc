CREATE TYPE status AS ENUM ('open', 'closed');
ALTER TYPE status ADD VALUE 'unknown';
ALTER TYPE status ADD VALUE IF NOT EXISTS 'unknown';