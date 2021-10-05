BEGIN;

CREATE TABLE IF NOT EXISTS "user" (
    id varchar(256) PRIMARY KEY,
    name varchar(256) NOT NULL
    );

CREATE TABLE IF NOT EXISTS user_tariff (
    user_id varchar(256) REFERENCES "user" (id) ON DELETE CASCADE,
    electricity_tariff money NOT NULL,
    cold_water_tariff money NOT NULL,
    hot_water_tariff money NOT NULL,
    out_water_tariff money NOT NULL,
    internet_tariff money
    );

CREATE TABLE IF NOT EXISTS user_meter_data (
    user_id varchar(256) REFERENCES "user" (id) ON DELETE CASCADE,
    electricity_meter integer NOT NULL,
    cold_water_meter integer NOT NULL,
    hot_water_meter integer NOT NULL
    );

COMMIT;
