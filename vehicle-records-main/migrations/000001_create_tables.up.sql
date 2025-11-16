CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
EXTENSION IF NOT EXISTS postgis;

-- Create table for vehicle records
CREATE TABLE vehicle_records
(
    -- Primary record identification
    record_id                               UUID PRIMARY KEY,
    record_store_time                       BIGINT         NOT NULL,
    record_send_time                        BIGINT         NOT NULL,

    -- Vehicle details
    citizen_vehicle_type                    INTEGER        NOT NULL,
    citizen_vehicle_color                   INTEGER        NOT NULL,
    citizen_vehicle_model                   TEXT,
    citizen_vehicle_distance                INTEGER        NOT NULL,
    citizen_vehicle_degree                  INTEGER        NOT NULL,
    citizen_plate_number                    TEXT           NOT NULL,
    citizen_vehicle_plate_number_type       INTEGER        NOT NULL,
    citizen_vehicle_plate_number_color      INTEGER        NOT NULL,


    -- Recognition and verification flags
    ocr_accuracy                            NUMERIC(6, 6)  NOT NULL,
    is_citizen_vehicle_distorted            BOOLEAN        NOT NULL,
    is_citizen_vehicle_plate_number_visible BOOLEAN        NOT NULL,
    citizen_park_type                       INTEGER        NOT NULL,

    -- Location and spatial references
    ring_id                                 BIGINT,
    street_id                               BIGINT,
    segment_id                              BIGINT,
    parking_lot_id                          BIGINT,


    tehran_request_id                       VARCHAR(250),
    plate_detection_id                      integer,
    retries                                 integer,

    -- User and system identifiers
    user_id                                 VARCHAR(250)    NOT NULL,
    lpr_vehicle_id                          VARCHAR(250)    NOT NULL,
    lpr_system_id                           TEXT           NOT NULL,
    lpr_system_app_id                       VARCHAR(250)    NOT NULL,
    lpr_system_app_version                  TEXT           NOT NULL,

    -- GPS and positioning data
    lpr_vehicle_gps_speed                   NUMERIC(10, 2) NOT NULL DEFAULT 0,
    lpr_vehicle_is_gps_signal_valid         BOOLEAN        NOT NULL,
    lpr_vehicle_gps_latitude                NUMERIC(10, 8) NOT NULL,
    lpr_vehicle_gps_longitude               NUMERIC(11, 8) NOT NULL,
    lpr_vehicle_gps_error                   INTEGER        NOT NULL,

    -- RTK (Real-Time Kinematic) positioning data
    lpr_vehicle_rtk_latitude                NUMERIC(10, 8) NOT NULL,
    lpr_vehicle_rtk_longitude               NUMERIC(11, 8) NOT NULL,
    lpr_vehicle_rtk_error                   INTEGER        NOT NULL,
    sent                                    BOOL                    DEFAULT false,
    created_at                              BIGINT         NOT NULL,
    updated_at                              BIGINT         NOT NULL,
    deleted_at                              BIGINT
);

-- Create index on frequently queried columns
CREATE INDEX idx_vehicle_records_user_id ON vehicle_records (user_id);
CREATE INDEX idx_vehicle_records_lpr_vehicle_id ON vehicle_records (lpr_vehicle_id);
CREATE INDEX idx_vehicle_records_record_store_time ON vehicle_records (record_store_time);


CREATE TABLE citizen_vehicle_photos
(
    id                  UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_id                          UUID          NOT NULL REFERENCES vehicle_records (record_id) ON DELETE CASCADE,
    photo_sequence_id                  INTEGER       NOT NULL,
    ocr_accuracy                       NUMERIC(6, 6) NOT NULL,
    lpr_vehicle_camera_id              INTEGER       NOT NULL,

    citizen_vehicle_photo              VARCHAR(355),
    citizen_vehicle_photo_area         TEXT,         -- Storing coordinates as text
    citizen_vehicle_plate_crop_photo   VARCHAR(355), -- Crop of plate photo

    citizen_vehicle_photo_capture_time BIGINT        NOT NULL,

    created_at                         BIGINT        NOT NULL,
    updated_at                         BIGINT        NOT NULL,
    deleted_at                         BIGINT
);

CREATE INDEX idx_citizen_vehicle_photos_record_id ON citizen_vehicle_photos (record_id);
CREATE INDEX idx_citizen_vehicle_photos_camera_id ON citizen_vehicle_photos (lpr_vehicle_camera_id);

CREATE TABLE states
(
    id                  UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id             VARCHAR(250),
    device_id           VARCHAR(250),
    creation_time       BIGINT           NOT NULL,

    -- Location data
    speed               FLOAT NULL,
    latitude            DOUBLE PRECISION NOT NULL,
    longitude           DOUBLE PRECISION NOT NULL,
    location_source     VARCHAR(20)      NOT NULL,

    -- Network data
    network_speed       INTEGER NULL,
    ip_address          INET             NOT NULL,
    mac_address         MACADDR          NOT NULL,
    adapter_name        VARCHAR(250)      NOT NULL,
    adapter_type        VARCHAR(250)      NOT NULL,
    adapter_description VARCHAR(255)     NOT NULL,

    created_at          BIGINT           NOT NULL,
    updated_at          BIGINT           NOT NULL,
    deleted_at          BIGINT
);