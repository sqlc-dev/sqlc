CREATE SCHEMA masterdata;

CREATE TYPE masterdata.rolling_stock_number_input AS (
    number           VARCHAR(20),
    number_type_code VARCHAR(20),
    rw_admin_id      BIGINT
);

CREATE FUNCTION masterdata.rolling_stock_get_or_create_ids(
    inputs   masterdata.rolling_stock_number_input[],
    on_date  TIMESTAMPTZ DEFAULT now(),
    train_id BIGINT DEFAULT NULL
)
    RETURNS TABLE
            (
                number           VARCHAR,
                rolling_stock_id BIGINT
            )
    LANGUAGE sql
    STABLE
AS
$$
SELECT '' :: VARCHAR AS number,
       0  :: BIGINT  AS rolling_stock_id
WHERE FALSE
$$;
