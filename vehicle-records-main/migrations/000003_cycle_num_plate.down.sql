ALTER TABLE records
    DROP COLUMN IF EXISTS cycle_id ,
    DROP COLUMN IF EXISTS citizen_plate_number_numeric ,
    DROP COLUMN IF EXISTS shamsi_time ;
