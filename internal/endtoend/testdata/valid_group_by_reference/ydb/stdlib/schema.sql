CREATE TABLE authors (
    id BigSerial,
    name Text NOT NULL,
    bio Text,
    PRIMARY KEY (id)
);

CREATE TABLE weather_metrics (
    time Timestamp NOT NULL,
    timezone_shift Int32,
    city_name Text,
    temp_c Double,
    feels_like_c Double,
    temp_min_c Double,
    temp_max_c Double,
    pressure_hpa Double,
    humidity_percent Double,
    wind_speed_ms Double,
    wind_deg Int32,
    rain_1h_mm Double,
    rain_3h_mm Double,
    snow_1h_mm Double,
    snow_3h_mm Double,
    clouds_percent Int32,
    weather_type_id Int32,
    PRIMARY KEY (time, city_name)
);

