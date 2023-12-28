CREATE TABLE IF NOT EXISTS job
(
    task_name text NOT NULL,
    last_run timestamp with time zone DEFAULT now() NOT NULL
);
