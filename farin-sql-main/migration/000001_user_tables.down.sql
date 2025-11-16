-- Drop tables in reverse order of creation to handle foreign key dependencies

DROP TABLE IF EXISTS rates;
DROP TABLE IF EXISTS exceptions;


DROP TABLE IF EXISTS parkings CASCADE;
DROP TABLE IF EXISTS segments CASCADE;
DROP TABLE IF EXISTS districts_parts CASCADE;
DROP TABLE IF EXISTS parts CASCADE;
DROP TABLE IF EXISTS districts_roads CASCADE;
DROP TABLE IF EXISTS roads CASCADE;
DROP TABLE IF EXISTS districts CASCADE;
DROP TABLE IF EXISTS cities CASCADE;
DROP TABLE IF EXISTS countries CASCADE;
DROP TABLE IF EXISTS contractors CASCADE;
-- Drop tables in reverse order of creation to handle foreign key dependencies

DROP TABLE IF EXISTS drivers_vehicles CASCADE;
DROP TABLE IF EXISTS vehicles CASCADE;
DROP TABLE IF EXISTS drivers CASCADE;
DROP TABLE IF EXISTS rings_districts CASCADE;
DROP TABLE IF EXISTS rings CASCADE;
DROP TABLE IF EXISTS contracts_districts CASCADE;
DROP TABLE IF EXISTS contractors_cities CASCADE;
DROP TABLE IF EXISTS contracts CASCADE;
DROP TABLE IF EXISTS contractors CASCADE;

-- Drop tables in reverse order of creation to handle foreign key dependencies

DROP TABLE IF EXISTS assignments CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS form_modules CASCADE;
DROP TABLE IF EXISTS modules CASCADE;
DROP TABLE IF EXISTS forms CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS peak_hour_multipliers CASCADE;
DROP TABLE IF EXISTS base_rates CASCADE;
DROP TABLE IF EXISTS road_categories CASCADE;
DROP TABLE IF EXISTS vehicle_categories CASCADE;
DROP TABLE IF EXISTS drivers_devices CASCADE;
DROP TABLE IF EXISTS driver_licenses CASCADE;
DROP TABLE IF EXISTS device_parts CASCADE;
DROP TABLE IF EXISTS devices CASCADE;

-- Drop ENUM types
DROP TYPE IF EXISTS gender CASCADE;
DROP TYPE IF EXISTS status CASCADE;
DROP TYPE IF EXISTS weekday CASCADE;
DROP TYPE IF EXISTS shift_work CASCADE;


