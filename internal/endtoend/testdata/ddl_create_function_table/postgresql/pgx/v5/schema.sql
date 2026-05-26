create type public.test_jsonb as (test_json jsonb);

create domain public.test_jsonb_domain as jsonb;

create table if not exists public.test_data
(
    test_id      integer,
    test_date    date,
    test_time    timestamp with time zone,
    test_string  text,
    test_varchar character varying,
    test_double  double precision,
    test_jsonb   test_jsonb,
    test_jsonb_domain   test_jsonb_domain
);

create function public.get_test(input_time timestamp without time zone DEFAULT now())
    returns TABLE
            (
                test_id      integer,
                test_date    date,
                test_time    timestamp with time zone,
                test_string  text,
                test_varchar character varying,
                test_double  double precision,
                test_jsonb  test_jsonb,
                test_jsonb_domain  test_jsonb_domain
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
       test_double,
       test_jsonb,
       test_jsonb_domain
FROM public.test_data
WHERE test_time <= input_time
$$;


CREATE FUNCTION public.get_all_tests_at_moment(p_target_time timestamp without time zone DEFAULT now())
    RETURNS TABLE
            (
                test_id                             integer,
                departure_date                       date,
                formation_time                       timestamp with time zone,
                disbanding_time                      timestamp with time zone,
                operation_id                         integer,
                operation                            text,
                operation_time                       timestamp with time zone,
                operation_test                    public.test_jsonb_domain,
                operation_test_code               text,
                operation_test_id                 integer,
                test_number                         character varying,
                departure_test                    public.test_jsonb_domain,
                departure_test_code               character varying,
                departure_test_id                 integer,
                index_number                         integer,
                destination_test                  public.test_jsonb_domain,
                destination_test_code             character varying,
                destination_test_id               integer,
                test_composition_wagons_count       integer,
                test_composition_net_weight         double precision,
                test_composition_gross_weight       double precision,
                test_composition_length             double precision,
                test_composition_conditional_length double precision,
                tests_count                          integer,
                tests_list                           text,
                testdriver_name                      text,
                testdriver_code                      text
            )
    LANGUAGE sql
    STABLE
AS
$$
SELECT 1                                                                         as test_id,
       CURRENT_DATE                                                              as departure_date,
       CURRENT_TIMESTAMP                                                         as formation_time,
       CURRENT_TIMESTAMP                                                         as disbanding_time,
       1                                                                         as operation_id,
       'Test Operation'                                                          as operation,
       CURRENT_TIMESTAMP                                                         as operation_time,
       ROW ('Test Station', 'Test Station')::public.test_jsonb_domain               as operation_test,
       'TST'                                                                     as operation_test_code,
       1                                                                         as operation_test_id,
       '123'                                                                     as test_number,
       ROW ('Departure Station', 'Departure Station')::public.test_jsonb_domain     as departure_test,
       'DEP'                                                                     as departure_test_code,
       1                                                                         as departure_test_id,
       1                                                                         as index_number,
       ROW ('Destination Station', 'Destination Station')::public.test_jsonb_domain as destination_test,
       'DST'                                                                     as destination_test_code,
       1                                                                         as destination_test_id,
       10                                                                        as test_composition_wagons_count,
       1000.50                                                                   as test_composition_net_weight,
       1500.75                                                                   as test_composition_gross_weight,
       250.0                                                                     as test_composition_length,
       200.0                                                                     as test_composition_conditional_length,
       2                                                                         as tests_count,
       'test1, test2'                                                            as tests_list,
       'John Doe'                                                                as testdriver_name,
       'JD001'                                                                   as testdriver_code
$$;
