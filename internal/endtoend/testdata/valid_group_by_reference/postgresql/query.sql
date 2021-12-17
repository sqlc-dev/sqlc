CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

-- name: ListAuthors :many
SELECT   id, name as name, bio
FROM     authors
GROUP BY name;

-- name: ListAuthorsIdenticalAlias :many
SELECT   id, name as name, bio
FROM     authors
GROUP BY name;


-- https://github.com/kyleconroy/sqlc/issues/1315

CREATE TABLE IF NOT EXISTS weather_metrics
(
    time             TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    timezone_shift   INT                         NULL,
    city_name        TEXT                        NULL,
    temp_c           DOUBLE PRECISION            NULL,
    feels_like_c     DOUBLE PRECISION            NULL,
    temp_min_c       DOUBLE PRECISION            NULL,
    temp_max_c       DOUBLE PRECISION            NULL,
    pressure_hpa     DOUBLE PRECISION            NULL,
    humidity_percent DOUBLE PRECISION            NULL,
    wind_speed_ms    DOUBLE PRECISION            NULL,
    wind_deg         INT                         NULL,
    rain_1h_mm       DOUBLE PRECISION            NULL,
    rain_3h_mm       DOUBLE PRECISION            NULL,
    snow_1h_mm       DOUBLE PRECISION            NULL,
    snow_3h_mm       DOUBLE PRECISION            NULL,
    clouds_percent   INT                         NULL,
    weather_type_id  INT                         NULL
);

-- name: ListMetrics :many
SELECT time_bucket('15 days', time) AS bucket, city_name, AVG(temp_c)
FROM weather_metrics
WHERE time > NOW() - (6 * INTERVAL '1 month')
GROUP BY bucket, city_name
ORDER BY bucket DESC;
