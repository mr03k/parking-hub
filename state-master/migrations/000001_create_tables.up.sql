CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE states
(
    id                     UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id                VARCHAR(250),
    record_id              varchar(100)             NOT NULL,
    lpr_system_id          VARCHAR(255)     NOT NULL,
    lpr_vehicle_id         varchar(100)             NOT NULL,
    lpr_system_app_id      varchar(100)             NOT NULL,
    lpr_system_app_version VARCHAR(20)      NOT NULL,
    lpr_vehicle_gps_latitude  DOUBLE PRECISION NOT NULL,
    lpr_vehicle_gps_longitude DOUBLE PRECISION NOT NULL,
    lpr_vehicle_gps_speed     FLOAT         NOT NULL,
    lpr_vehicle_gps_error     INTEGER       NOT NULL,
    record_store_time      BIGINT          NOT NULL,
    record_send_time       BIGINT          NOT NULL,
    server_availability    BOOLEAN         NOT NULL,
    server_ping_time       INTEGER         NOT NULL,

    created_at             BIGINT          NOT NULL,
    updated_at             BIGINT          NOT NULL,
    deleted_at             BIGINT
);

-- Create index for record_id for faster lookups
CREATE INDEX idx_states_record_id ON states(record_id);

-- Create index for record_store_time for time-based queries
CREATE INDEX idx_states_record_store_time ON states(record_store_time);

-- Create index for record_send_time for time-based queries
CREATE INDEX idx_states_record_send_time ON states(record_send_time);