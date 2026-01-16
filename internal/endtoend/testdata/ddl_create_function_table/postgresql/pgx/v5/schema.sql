-- Create mock test table
create table if not exists public.test_data
(
    test_id      integer,
    test_date    date,
    test_time    timestamp with time zone,
    test_string  text,
    test_varchar character varying,
    test_double  double precision
);

create function public.get_test(input_time timestamp without time zone DEFAULT now())
    returns TABLE
            (
                test_id      integer,
                test_date    date,
                test_time    timestamp with time zone,
                test_string  text,
                test_varchar character varying,
                test_double  double precision
            )
    stable
    language sql
as
$$
SELECT test_id,
       test_date,
       test_time,
       test_string,
       test_varchar,
       test_double
FROM public.test_data
WHERE test_time <= input_time
$$;
