CREATE TABLE todos (
    id   serial PRIMARY KEY,
    task text
);

CREATE PROCEDURE create_todo(IN p_task text, OUT p_id int)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO todos (task) VALUES (p_task) RETURNING id INTO p_id;
END;
$$;
