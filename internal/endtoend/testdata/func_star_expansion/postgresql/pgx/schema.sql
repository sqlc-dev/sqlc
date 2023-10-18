create function test_func_get_time ()
    returns timestamp AS $$
Begin
    return now();
End;
$$ language plpgsql;

create function test_func_select_string (in p_string text)
    returns text AS $$
Begin
    return p_string;
End;
$$ language plpgsql;

create function test_func_select_blog(in p_id int)
    returns table (id int, name text, created_at timestamp, updated_at timestamp) AS $$
BEGIN RETURN QUERY
    select id, name, created_at, updated_at from blog where id = p_id;
END;
$$ language plpgsql;
