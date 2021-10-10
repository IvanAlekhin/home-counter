BEGIN;

CREATE TABLE IF NOT EXISTS "user" (
    id varchar(256) PRIMARY KEY,
    name varchar(256) NOT NULL
    );

CREATE TABLE IF NOT EXISTS user_tariff (
    user_id varchar(256) REFERENCES "user" (id) ON DELETE CASCADE,
    electricity money NOT NULL,
    cold_water money NOT NULL,
    hot_water money NOT NULL,
    out_water money NOT NULL,
    internet money
    );

CREATE TABLE IF NOT EXISTS user_meter_data (
    user_id varchar(256) REFERENCES "user" (id) ON DELETE CASCADE,
    electricity integer NOT NULL,
    cold_water integer NOT NULL,
    hot_water integer NOT NULL
    );

COMMIT;
