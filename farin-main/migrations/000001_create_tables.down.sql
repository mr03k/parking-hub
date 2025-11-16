DROP TABLE IF EXISTS driver_assignments;
DROP TABLE IF EXISTS lpr_vehicles_combination CASCADE;

DROP TABLE IF EXISTS vehicles CASCADE;

DROP TABLE IF EXISTS license_plate_reader_devices CASCADE;

DROP TABLE IF EXISTS contractors CASCADE;

DROP TABLE IF EXISTS contracts CASCADE;

|DROP TABLE IF EXISTS users CASCADE;

DROP TABLE IF EXISTS drivers CASCADE;

DROP TABLE IF EXISTS calendars;
DROP TABLE IF EXISTS rings;
DROP TABLE IF EXISTS vehicle_records;
DROP TABLE IF EXISTS citizen_vehicle_photos;

DROP INDEX IF EXISTS idx_device_locations_coords;
DROP INDEX IF EXISTS idx_device_locations_creation_time;
DROP INDEX IF EXISTS idx_device_locations_device_id;
DROP TABLE IF EXISTS device_locations;


DROP EXTENSION IF EXISTS "uuid-ossp";

