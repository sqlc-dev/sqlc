create table blog (
    id serial primary key,
    name text not null
);

create function test_select_blog(in p_id int)
    returns table (id int, name text) AS $$
BEGIN RETURN QUERY
    select id, name from blog where id = p_id;
END;
$$ language plpgsql;
