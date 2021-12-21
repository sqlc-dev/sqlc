CREATE TABLE authors (
  id   BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name text    NOT NULL,
  bio  text,
  UNIQUE(name)
);

-- name: ListAuthors :many
SELECT   id, name as full_name, bio
FROM     authors
GROUP BY full_name;

-- name: ListAuthorsIdenticalAlias :many
SELECT   id, name as name, bio
FROM     authors
GROUP BY name;


-- https://github.com/kyleconroy/sqlc/issues/1315

CREATE TABLE IF NOT EXISTS weather_metrics
(
    time             TIMESTAMP NOT NULL,
    timezone_shift   INT       NULL,
    city_name        TEXT      NULL,
    temp_c           FLOAT     NULL,
    feels_like_c     FLOAT     NULL,
    temp_min_c       FLOAT     NULL,
    temp_max_c       FLOAT     NULL,
    pressure_hpa     FLOAT     NULL,
    humidity_percent FLOAT     NULL,
    wind_speed_ms    FLOAT     NULL,
    wind_deg         INT       NULL,
    rain_1h_mm       FLOAT     NULL,
    rain_3h_mm       FLOAT     NULL,
    snow_1h_mm       FLOAT     NULL,
    snow_3h_mm       FLOAT     NULL,
    clouds_percent   INT       NULL,
    weather_type_id  INT       NULL
);

-- name: ListMetrics :many
SELECT time_bucket('15 days', time) AS bucket, city_name, AVG(temp_c)
FROM weather_metrics
WHERE DATE_SUB(NOW(), INTERVAL 6 MONTH)
GROUP BY bucket, city_name
ORDER BY bucket DESC;
