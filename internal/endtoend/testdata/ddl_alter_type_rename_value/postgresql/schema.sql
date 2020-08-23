CREATE TYPE status AS ENUM ('open', 'closed');
ALTER TYPE status RENAME VALUE 'closed' TO 'shut';