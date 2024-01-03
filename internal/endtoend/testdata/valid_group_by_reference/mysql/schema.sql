CREATE TABLE authors (
  id   BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(10)    NOT NULL,
  bio  text,
  UNIQUE(name)
);

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
