-- name: DoStuff :exec
DO $$
    BEGIN
        ALTER TABLE authors
        ADD COLUMN marked_for_processing bool;
    END
$$;
