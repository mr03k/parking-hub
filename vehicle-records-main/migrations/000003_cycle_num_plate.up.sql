ALTER TABLE vehicle_records
    ADD COLUMN cycle_id INTEGER ,
    ADD COLUMN citizen_plate_number_numeric      BIGINT,
    ADD COLUMN shamsi_time VARCHAR(60);