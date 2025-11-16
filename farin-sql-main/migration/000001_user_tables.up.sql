CREATE
EXTENSION postgis;

CREATE TABLE countries
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each country
    country_name VARCHAR(100) NOT NULL,                      -- Name of the country
    country_code VARCHAR(3) UNIQUE,                          -- ISO 3166-1 alpha-3 code (e.g., 'USA' for the United States)
    iso_code     VARCHAR(2) UNIQUE,                          -- ISO 3166-1 alpha-2 code (e.g., 'US' for the United States)
    region       VARCHAR(50),                                -- Geographic region (e.g., Asia, Europe, Africa)
    capital      VARCHAR(100),                               -- Capital city of the country
    phone_code   VARCHAR(10),                                -- International phone code (e.g., '+98' for Iran)
    currency     VARCHAR(50),                                -- Currency used in the country (e.g., Rial, Dollar, Euro)
    population   BIGINT,                                     -- Population of the country (optional)
    area         FLOAT,                                      -- Area of the country in square kilometers (optional)
    geo_boundary GEOMETRY(POLYGON),                          -- Geographical boundary of the country (optional, as a polygon)
    created_at   INT          NOT NULL
);

CREATE TABLE cities
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each city
    city_name  VARCHAR(100) NOT NULL,                      -- Name of the city
    code_city  VARCHAR(3) UNIQUE,                          -- Unique code for each city
    id_country UUID         NOT NULL,                      -- Foreign key related to the countries table
    boundary   GEOGRAPHY(POLYGON),                         -- Optional geographical boundary of the city as a polygon
    created_at INT          NOT NULL,
    CONSTRAINT fk_country
        FOREIGN KEY (id_country)                           -- Foreign key constraint
            REFERENCES countries (id)                      -- References the 'id' field in the 'countries' table
            ON DELETE CASCADE                              -- Cascade delete if the country is deleted
);

CREATE TABLE districts
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each district
    district_name VARCHAR(100) NOT NULL,                      -- Name of the district
    code_district VARCHAR(10) UNIQUE,                         -- Unique code for each district
    id_city       UUID         NOT NULL,                      -- Foreign key related to the cities table
    boundary_geo  GEOGRAPHY(POLYGON),                         -- Optional geographical boundary of the district as a polygon
    population    BIGINT,                                     -- Optional population of the district
    area          FLOAT,                                      -- Optional area of the district in square kilometers
    created_at    INT          NOT NULL,
    CONSTRAINT fk_city
        FOREIGN KEY (id_city)                                 -- Foreign key constraint linking to cities table
            REFERENCES cities (id)                            -- Reference to 'id' field in the 'cities' table
            ON DELETE CASCADE                                 -- Cascade delete if the city is deleted
);

CREATE TABLE roads
(
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),              -- Unique identifier for each road
    road_name              VARCHAR(100)       NOT NULL,                             -- Name of the road
    code_road              VARCHAR(10) UNIQUE NOT NULL,                             -- Unique code to identify each road
    type_road              VARCHAR(50),                                             -- Type of road (e.g., main road, highway)
    grade_road             VARCHAR(1) CHECK (grade_road IN ('الف', 'ب', 'ج', 'د')), -- Road grade (as an ENUM)
    length_road            DECIMAL(10, 2),                                          -- Length of the road in kilometers
    width_road             DECIMAL(5, 2),                                           -- Width of the road in meters
    limit_speed            INT,                                                     -- Maximum speed limit in kilometers per hour
    boundary_road          GEOGRAPHY(POLYGON),                                      -- Geographical boundary of the road (optional, as a polygon)
    spots_parking          INT,                                                     -- Number of parking spots on the road
    spots_parking_disabled INT,                                                     -- Number of disabled parking spots on the road
    created_at             INT                NOT NULL,
    description            TEXT                                                     -- Additional description or notes about the road
);


CREATE TABLE districts_roads
(
    id_road     UUID NOT NULL,         -- Foreign key related to the roads table
    id_district UUID NOT NULL,         -- Foreign key related to the districts table
    CONSTRAINT fk_road
        FOREIGN KEY (id_road)          -- Foreign key constraint for roads
            REFERENCES roads (id)      -- Reference to 'id' in roads table
            ON DELETE CASCADE,         -- Cascade delete when a road is deleted
    CONSTRAINT fk_district
        FOREIGN KEY (id_district)      -- Foreign key constraint for districts
            REFERENCES districts (id)  -- Reference to 'id' in districts table
            ON DELETE CASCADE,         -- Cascade delete when a district is deleted
    PRIMARY KEY (id_road, id_district) -- Composite primary key to ensure uniqueness of the relation
);

CREATE TABLE parts
(
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each part (segment)
    part_name              VARCHAR(100)       NOT NULL,                -- Name of the part (segment)
    code_part              VARCHAR(10) UNIQUE NOT NULL,                -- Unique code to identify each part (segment)
    id_road                UUID               NOT NULL,                -- Foreign key related to the roads table
    length_part            DECIMAL(10, 2),                             -- Length of the part in kilometers
    boundary_part          GEOGRAPHY(POLYGON),                         -- Geographical boundary of the part as a polygon
    spots_parking          INT,                                        -- Number of parking spots in this part
    spots_parking_disabled INT,                                        -- Number of disabled parking spots in this part
    description            TEXT,                                       -- Additional description or notes about the part
    CONSTRAINT fk_road
        FOREIGN KEY (id_road)                                          -- Foreign key constraint for roads
            REFERENCES roads (id)                                      -- References the 'id' field in the 'roads' table
            ON DELETE CASCADE                                          -- Cascade delete if the road is deleted
);


CREATE TABLE districts_parts
(
    id_part     UUID NOT NULL,         -- Foreign key related to the parts (road segments) table
    id_district UUID NOT NULL,         -- Foreign key related to the districts table
    CONSTRAINT fk_part
        FOREIGN KEY (id_part)          -- Foreign key constraint for parts
            REFERENCES parts (id)      -- Reference to 'id' in parts table
            ON DELETE CASCADE,         -- Cascade delete when a part is deleted
    CONSTRAINT fk_district
        FOREIGN KEY (id_district)      -- Foreign key constraint for districts
            REFERENCES districts (id)  -- Reference to 'id' in districts table
            ON DELETE CASCADE,         -- Cascade delete when a district is deleted
    PRIMARY KEY (id_part, id_district) -- Composite primary key to ensure uniqueness of the relation
);

CREATE TABLE segments
(
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each segment
    segment_name           VARCHAR(100)       NOT NULL,                -- Name of the segment
    segment_code           VARCHAR(10) UNIQUE NOT NULL,                -- Unique code to identify each segment
    part_id                UUID               NOT NULL,                -- Foreign key related to the parts table
    road_id                UUID               NOT NULL,                -- Foreign key related to the roads table
    district_id            UUID               NOT NULL,                -- Foreign key related to the districts table
    segment_length         DECIMAL(10, 2),                             -- Length of the segment in kilometers
    segment_boundary       GEOMETRY(POLYGON),                          -- Geographic boundary of the segment as a polygon
    parking_spots          INT              DEFAULT 0,                 -- Number of parking spots in the segment
    disabled_parking_spots INT              DEFAULT 0,                 -- Number of parking spots for disabled people
    description            TEXT,                                       -- Additional description
    created_at             INT                NOT NULL,                -- Timestamp when the record is created

    -- Foreign key constraints
    CONSTRAINT fk_part
        FOREIGN KEY (part_id)
            REFERENCES parts (id)
            ON DELETE CASCADE,                                         -- Cascade delete when a part is deleted

    CONSTRAINT fk_road
        FOREIGN KEY (road_id)
            REFERENCES roads (id)
            ON DELETE CASCADE,                                         -- Cascade delete when a road is deleted

    CONSTRAINT fk_district
        FOREIGN KEY (district_id)
            REFERENCES districts (id)
            ON DELETE CASCADE                                          -- Cascade delete when a district is deleted
);


CREATE TABLE parkings
(
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),                                -- Unique identifier for each parking spot
    code_parking        VARCHAR(50) UNIQUE NOT NULL,                                               -- Unique code to identify each parking spot
    id_segment          UUID               NOT NULL,                                               -- Foreign key related to the segments table
    type_parking        VARCHAR(10)        NOT NULL,                                               -- Type of parking (e.g., regular or disabled)
    boundary_parking    GEOGRAPHY(POLYGON),                                                        -- Geographical boundary of the parking spot as a polygon
    position            VARCHAR(1) CHECK (position IN ('R', 'L')),                                 -- Position of the parking (right or left of the segment)
    status_availability VARCHAR(20) CHECK (status_availability IN ('Full', 'Empty', 'Available')), -- Status of availability
    description         TEXT,                                                                      -- Additional description or notes about the parking spot
    created_at          INT                NOT NULL,
    CONSTRAINT fk_segment
        FOREIGN KEY (id_segment)                                                                   -- Foreign key constraint for segments
            REFERENCES segments (id)                                                               -- Reference to 'id' in the segments table
            ON DELETE CASCADE                                                                      -- Cascade delete when a segment is deleted
);

CREATE TABLE contractors
(
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each contractor
    contractor_name        VARCHAR(100)       NOT NULL,                -- Name of the contractor
    code_contractor        VARCHAR(10) UNIQUE NOT NULL,                -- Unique code for identifying the contractor
    number_registration    VARCHAR(50),                                -- Registration number of the contractor (for companies)
    person_contact         VARCHAR(100),                               -- Name of the responsible person or representative
    ceo_name               VARCHAR(100),                               -- Name of the CEO
    signatories_authorized TEXT,                                       -- Names of authorized signatories
    phone_number           VARCHAR(15),                                -- Contact phone number
    email                  VARCHAR(100),                               -- Email address
    address                VARCHAR(255),                               -- Office address of the contractor
    type_contract          VARCHAR(50),                                -- Type of contract (e.g., equipment and operation, only operation)
    number_account_bank    VARCHAR(30),                                -- Bank account number for payments
    created_at             INT                NOT NULL,
    description            TEXT                                        -- Additional description or notes
);


CREATE TABLE contracts
(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each contract
    number_contract  VARCHAR(50) UNIQUE NOT NULL,                -- Unique contract number
    date_contract    DATE               NOT NULL,                -- Date the contract was issued
    date_start       DATE               NOT NULL,                -- Start date of the contract
    date_end         DATE               NOT NULL,                -- End date of the contract
    amount_contract  BIGINT             NOT NULL,                -- Total amount of the contract
    type_contract    VARCHAR(50)        NOT NULL,                -- Type of contract (e.g., equipment and operation, only operation)
    id_contractor    UUID               NOT NULL,                -- Foreign key related to the contractors table
    period_operation INT,                                        -- Duration of parking operation in days
    period_equipment INT,                                        -- Duration of parking equipment setup in days
    description      TEXT,                                       -- Additional description or special conditions related to the contract
    created_at       INT                NOT NULL,
    CONSTRAINT fk_contractor
        FOREIGN KEY (id_contractor)                              -- Foreign key constraint for contractors
            REFERENCES contractors (id)                          -- Reference to 'id' in contractors table
            ON DELETE CASCADE                                    -- Cascade delete when a contractor is deleted
);


CREATE TABLE contractors_cities
(
    id_contract UUID NOT NULL,         -- Foreign key related to the contracts table
    id_city     UUID NOT NULL,         -- Foreign key related to the cities table
    CONSTRAINT fk_contract
        FOREIGN KEY (id_contract)      -- Foreign key constraint for contracts
            REFERENCES contracts (id)  -- Reference to 'id' in the contracts table
            ON DELETE CASCADE,         -- Cascade delete when a contract is deleted
    CONSTRAINT fk_city
        FOREIGN KEY (id_city)          -- Foreign key constraint for cities
            REFERENCES cities (id)     -- Reference to 'id' in the cities table
            ON DELETE CASCADE,         -- Cascade delete when a city is deleted
    PRIMARY KEY (id_contract, id_city) -- Composite primary key to ensure uniqueness of the relation
);


CREATE TABLE contracts_districts
(
    id_contract UUID NOT NULL,             -- Foreign key related to the contracts table
    id_district UUID NOT NULL,             -- Foreign key related to the districts table
    CONSTRAINT fk_contract
        FOREIGN KEY (id_contract)          -- Foreign key constraint for contracts
            REFERENCES contracts (id)      -- Reference to 'id' in the contracts table
            ON DELETE CASCADE,             -- Cascade delete when a contract is deleted
    CONSTRAINT fk_district
        FOREIGN KEY (id_district)          -- Foreign key constraint for districts
            REFERENCES districts (id)      -- Reference to 'id' in the districts table
            ON DELETE CASCADE,             -- Cascade delete when a district is deleted
    PRIMARY KEY (id_contract, id_district) -- Composite primary key to ensure uniqueness of the relation
);

CREATE TABLE rings
(
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each ring
    name_ring              VARCHAR(100)       NOT NULL,                -- Name or number of the ring
    code_ring              VARCHAR(10) UNIQUE NOT NULL,                -- Unique code to identify each ring
    length_ring            DECIMAL(10, 2)     NOT NULL,                -- Length of the ring in kilometers
    boundary_ring          GEOGRAPHY(POLYGON),                         -- Geographical boundary of the ring as a polygon
    spots_parking          INT,                                        -- Number of regular parking spots in the ring
    spots_parking_disabled INT,                                        -- Number of disabled parking spots in the ring
    signs_traffic          INT,                                        -- Number of regular traffic signs in the ring
    signs_traffic_disabled INT,                                        -- Number of disabled traffic signs in the ring
    point_start            GEOGRAPHY(POINT),                           -- Starting point of the ring as a geographic point
    distance_buffer        DECIMAL(5, 2),                              -- Buffer distance related to the starting point in meters
    created_at             INT                NOT NULL,
    description            TEXT                                        -- Additional description or notes about the ring
);

CREATE TABLE rings_districts
(
    id_ring     UUID NOT NULL,         -- Foreign key related to the rings table
    id_district UUID NOT NULL,         -- Foreign key related to the districts table
    CONSTRAINT fk_ring
        FOREIGN KEY (id_ring)          -- Foreign key constraint for rings
            REFERENCES rings (id)      -- Reference to 'id' in the rings table
            ON DELETE CASCADE,         -- Cascade delete when a ring is deleted
    CONSTRAINT fk_district
        FOREIGN KEY (id_district)      -- Foreign key constraint for districts
            REFERENCES districts (id)  -- Reference to 'id' in the districts table
            ON DELETE CASCADE,         -- Cascade delete when a district is deleted
    PRIMARY KEY (id_ring, id_district) -- Composite primary key to ensure uniqueness of the relation
);

CREATE TABLE drivers
(
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each driver
    first_name                  VARCHAR(50)        NOT NULL,                -- First name of the driver
    name_last                   VARCHAR(50)        NOT NULL,                -- Last name of the driver
    gender                      VARCHAR(10),                                -- Gender of the driver (e.g., male, female)
    code_driver                 VARCHAR(10) UNIQUE NOT NULL,                -- Unique code to identify the driver
    id_national                 VARCHAR(20),                                -- National ID number of the driver
    code_postal                 VARCHAR(10),                                -- Postal code of the driver's residence
    number_phone                VARCHAR(15),                                -- Fixed phone number of the driver
    number_mobile               VARCHAR(15),                                -- Mobile phone number of the driver
    email                       VARCHAR(100),                               -- Email address of the driver
    address                     TEXT,                                       -- Residence address of the driver
    id_contractor               UUID,                                       -- Foreign key related to the contractors table
    type_driver                 VARCHAR(10),                                -- Type of driver (e.g., main, reserve)
    type_shift                  VARCHAR(10),                                -- Type of shift (e.g., morning, evening, both)
    status_employment           VARCHAR(20),                                -- Employment status of the driver (e.g., active, inactive)
    date_start_employment       DATE,                                       -- Start date of employment
    date_end_employment         DATE,                                       -- End date of employment (if applicable)
    driver_photo                VARCHAR(200),                               -- Photo of the driver
    image_card_id               VARCHAR(200),                               -- Image of the national ID card
    birth_certificate_image     VARCHAR(200),                               -- Image of the birth certificate
    image_card_service_military VARCHAR(200),                               -- Image of the military service completion card
    image_certificate_health    VARCHAR(200),                               -- Image of the health certificate
    image_record_criminal       VARCHAR(200),                               -- Image of the criminal record certificate
    created_at                  INT                NOT NULL,
    description                 TEXT                                        -- Additional description or notes
);

CREATE TABLE vehicles
(
    id                           UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each vehicle
    code_vehicle                 VARCHAR(20) UNIQUE NOT NULL,                -- Unique code to identify each vehicle
    vin                          VARCHAR(20) UNIQUE NOT NULL,                -- Vehicle Identification Number (VIN)
    plate_license                VARCHAR(15) UNIQUE NOT NULL,                -- License plate number
    type_vehicle                 VARCHAR(50),                                -- Type of vehicle (e.g., motorcycle, sedan, minibus, truck)
    brand                        VARCHAR(50),                                -- Brand of the vehicle (e.g., Peugeot, Hyundai)
    model                        VARCHAR(50),                                -- Model of the vehicle (e.g., 206, Sonata)
    color                        VARCHAR(30),                                -- Color of the vehicle
    manufacture_of_year          INT,                                        -- Year of manufacture
    kilometers_initial           BIGINT,                                     -- Initial kilometers at start of operation
    expiry_insurance_party_third DATE,                                       -- Expiry date of third-party insurance
    expiry_insurance_body        DATE,                                       -- Expiry date of body insurance
    image_document_vehicle       VARCHAR(200),                               -- Image of the vehicle's document
    image_card_vehicle           VARCHAR(200),                               -- Image of the vehicle's card
    third_party_insurance_image  VARCHAR(200),                               -- Image of third-party insurance
    body_insurance_image         VARCHAR(200),                               -- Image of body insurance
    id_contractor                UUID,                                       -- Foreign key related to the contractors table
    status                       VARCHAR(20),                                -- Status of the vehicle (e.g., active, inactive, under repair)
    description                  TEXT,                                       -- Additional description or notes
    created_at                   INT                NOT NULL,
    CONSTRAINT fk_contractor
        FOREIGN KEY (id_contractor)                                          -- Foreign key constraint for contractors
            REFERENCES contractors (id)                                      -- Reference to 'id' in the contractors table
            ON DELETE SET NULL                                               -- Set to NULL if the related contractor is deleted
);


CREATE TABLE drivers_vehicles
(
    id_driver UUID NOT NULL,         -- Foreign key related to the drivers table
    id_veh    UUID NOT NULL,         -- Foreign key related to the vehicles table
    CONSTRAINT fk_driver
        FOREIGN KEY (id_driver)      -- Foreign key constraint for drivers
            REFERENCES drivers (id)  -- Reference to 'id' in the drivers table
            ON DELETE CASCADE,       -- Cascade delete when a driver is deleted
    CONSTRAINT fk_vehicle
        FOREIGN KEY (id_veh)         -- Foreign key constraint for vehicles
            REFERENCES vehicles (id) -- Reference to 'id' in the vehicles table
            ON DELETE CASCADE,       -- Cascade delete when a vehicle is deleted
    PRIMARY KEY (id_driver, id_veh)  -- Composite primary key to ensure uniqueness of the relation
);

CREATE TABLE devices
(
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each device
    code_device           VARCHAR(20) UNIQUE NOT NULL,                -- Unique code to identify each device
    number_serial         VARCHAR(50) UNIQUE NOT NULL,                -- Serial number of the device
    model                 VARCHAR(50),                                -- Model of the device
    date_installation     DATE,                                       -- Date of installation on the vehicle
    date_expiry_warranty  DATE,                                       -- Expiry date of the device warranty
    date_expiry_insurance DATE,                                       -- Expiry date of the device insurance
    class_device          VARCHAR(50),                                -- Class of the device (e.g., fixed, mobile)
    image_contract        VARCHAR(200),                               -- Image of the device's contract
    image_insurance       VARCHAR(200),                               -- Image of the device's insurance
    id_contractor         UUID,                                       -- Foreign key related to the contractors table
    description           TEXT,                                       -- Additional description or notes
    created_at            INT                NOT NULL,
    CONSTRAINT fk_contractor
        FOREIGN KEY (id_contractor)                                   -- Foreign key constraint for contractors
            REFERENCES contractors (id)                               -- Reference to 'id' in the contractors table
            ON DELETE SET NULL                                        -- Set to NULL if the related contractor is deleted
);
CREATE TABLE device_parts
(
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each part
    code_part            VARCHAR(20) UNIQUE NOT NULL,                -- Unique code to identify each part
    part_name            VARCHAR(100),                               -- Name of the part
    type_part            VARCHAR(50),                                -- Type of the part (e.g., camera, lens, microphone)
    brand                VARCHAR(50),                                -- Brand of the part
    model                VARCHAR(50),                                -- Model of the part
    number_serial        VARCHAR(50),                                -- Serial number of the part
    date_installation    DATE,                                       -- Date of installation
    date_expiry_warranty DATE,                                       -- Expiry date of the warranty
    period_maintenance   VARCHAR(50),                                -- Maintenance and repair period (e.g., 6 months)
    id_device            UUID,                                       -- Foreign key related to the devices table
    description          TEXT,                                       -- Additional description or notes
    created_at           INT                NOT NULL,
    CONSTRAINT fk_device
        FOREIGN KEY (id_device)                                      -- Foreign key constraint for devices
            REFERENCES devices (id)                                  -- Reference to 'id' in the devices table
            ON DELETE SET NULL                                       -- Set to NULL if the related device is deleted
);

CREATE TABLE driver_licenses
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each license
    id_driver      UUID               NOT NULL,                -- Foreign key related to the drivers table
    license_number VARCHAR(20) UNIQUE NOT NULL,                -- Unique license number
    type_license   VARCHAR(50),                                -- Type of the license (e.g., Class 1, Class 2, Motorcycle)
    date_issue     DATE,                                       -- Date of issuance
    date_expiry    DATE,                                       -- Expiry date of the license
    image_license  VARCHAR(200),                               -- Image of the license
    description    TEXT,                                       -- Additional description or notes
    created_at     INT                NOT NULL,
    CONSTRAINT fk_driver
        FOREIGN KEY (id_driver)                                -- Foreign key constraint for drivers
            REFERENCES drivers (id)                            -- Reference to 'id' in the drivers table
            ON DELETE CASCADE                                  -- Cascade delete if the related driver is deleted
);


CREATE TABLE drivers_devices
(
    id_driver     UUID NOT NULL,           -- Foreign key related to the drivers table
    id_device_lpr UUID NOT NULL,           -- Foreign key related to the license plate recognition devices table
    CONSTRAINT fk_driver
        FOREIGN KEY (id_driver)            -- Foreign key constraint for drivers
            REFERENCES drivers (id)        -- Reference to 'id' in the drivers table
            ON DELETE CASCADE,             -- Cascade delete when a driver is deleted
    CONSTRAINT fk_device_lpr
        FOREIGN KEY (id_device_lpr)        -- Foreign key constraint for devices
            REFERENCES devices (id)        -- Reference to 'id' in the devices table
            ON DELETE CASCADE,             -- Cascade delete when a device is deleted
    PRIMARY KEY (id_driver, id_device_lpr) -- Composite primary key to ensure uniqueness of the relation
);

CREATE TABLE vehicle_categories
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each category
    code_category VARCHAR(20) UNIQUE NOT NULL,                -- Unique code to identify each category
    name_category VARCHAR(50)        NOT NULL,                -- Name of the category (e.g., Passenger and Pickup, Truck and Minibus, Motorcycle)
    description   TEXT,                                       -- Additional description or notes about the category
    created_at    INT                NOT NULL
);

CREATE TABLE road_categories
(
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each road category
    code_category_road VARCHAR(10) UNIQUE NOT NULL,                -- Unique code to identify each road category
    name_category_road VARCHAR(50)        NOT NULL,                -- Name of the road category (e.g., Main Road, Secondary Road, Highway)
    description        TEXT,                                       -- Additional description or notes about the road category
    created_at         INT                NOT NULL
);


CREATE TABLE base_rates
(
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each base rate
    id_category_vehicle UUID           NOT NULL,                    -- Foreign key related to the vehicle categories table
    from_minutes        INT            NOT NULL,                    -- Starting minute for the base rate
    to_minutes          INT            NOT NULL,                    -- Ending minute for the base rate
    base_rate           DECIMAL(10, 2) NOT NULL,                    -- Base rate amount
    description         TEXT,                                       -- Additional description or notes about the rate
    created_at          INT            NOT NULL,
    CONSTRAINT fk_category_vehicle
        FOREIGN KEY (id_category_vehicle)                           -- Foreign key constraint for vehicle categories
            REFERENCES vehicle_categories (id)                      -- Reference to 'id' in the vehicle_categories table
            ON DELETE CASCADE                                       -- Cascade delete if the related vehicle category is deleted
);


CREATE TABLE peak_hour_multipliers
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each peak hour multiplier
    code_time_peak VARCHAR(20) UNIQUE NOT NULL,                -- Unique code to identify each peak hour period
    description    VARCHAR(100),                               -- Description of the peak hour (e.g., "Morning" or "Evening")
    multiplier     DECIMAL(10, 2)     NOT NULL,                -- Multiplier to be added to the rates during the peak hour
    weekday        VARCHAR(10),                                -- Day of the week (e.g., "Saturday", "Sunday")
    time_start     varchar(20)        NOT NULL,                -- Start time of the peak hour
    time_end       varchar(20)        NOT NULL,                -- End time of the peak hour
    from_valid     DATE               NOT NULL,                -- Start date of validity for the multiplier
    to_valid       DATE               NOT NULL,                -- End date of validity for the multiplier
    flag           VARCHAR(20),                                -- Status (e.g., "Year-Round" or "Exception")
    created_at     INT                NOT NULL
);

CREATE TYPE gender AS ENUM ('Male', 'Female', 'Other');
CREATE TYPE status AS ENUM ('Active', 'Inactive');

CREATE TABLE users
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each peak hour multiplier
    username        VARCHAR(50) UNIQUE NOT NULL,
    password        VARCHAR(255)       NOT NULL,
    first_name      VARCHAR(100)       NOT NULL,
    last_name       VARCHAR(100)       NOT NULL,
    email           VARCHAR(100)       NOT NULL,
    number_phone    VARCHAR(15),
    number_mobile   VARCHAR(15),
    id_national     VARCHAR(10),
    code_postal     VARCHAR(10),
    name_company    VARCHAR(100),
    image_profile   VARCHAR(200),
    msisdn          VARCHAR(255)       NOT NULL UNIQUE,
    msisdn_verified BOOLEAN          DEFAULT TRUE,
    gender          gender,                                     -- Use the defined ENUM type
    address         VARCHAR(255),
    status          status             NOT NULL,                -- Use the defined ENUM type
    created_at      INT                NOT NULL
);


CREATE TABLE forms
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each peak hour multiplier
    form_name   VARCHAR(100) NOT NULL,                      -- Name of the form
    description VARCHAR(255),                               -- Description related to the form
    created_at  INT          NOT NULL
);


CREATE TABLE modules
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each peak hour multiplier
    module_name VARCHAR(100) NOT NULL,                      -- Name of the module
    description VARCHAR(255),                               -- Description related to the module
    created_at  INT          NOT NULL
);



CREATE TABLE form_modules
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each peak hour multiplier
    id_module  UUID NOT NULL,                              -- Foreign key to the modules table
    id_form    UUID NOT NULL,                              -- Foreign key to the forms table
    created_at INT  NOT NULL,
    CONSTRAINT fk_module
        FOREIGN KEY (id_module)                            -- Foreign key constraint for modules
            REFERENCES modules (id)                        -- Reference to 'id_module' in the modules table
            ON DELETE CASCADE,                             -- Cascade delete if the related module is deleted
    CONSTRAINT fk_form
        FOREIGN KEY (id_form)                              -- Foreign key constraint for forms
            REFERENCES forms (id)                          -- Reference to 'id_form' in the forms table
            ON DELETE CASCADE                              -- Cascade delete if the related form is deleted
);


CREATE TABLE roles
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name   VARCHAR(50) NOT NULL, -- Name of the role
    description VARCHAR(255),         -- Additional description about the role
    created_at  INT         NOT NULL  -- Timestamp for when the role was created
);


CREATE TABLE assignments
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_user         UUID NOT NULL,               -- Foreign key to the users table
    id_role         UUID NOT NULL,               -- Foreign key to the roles table
    id_module       UUID NOT NULL,               -- Foreign key to the modules table
    id_form         UUID NOT NULL,               -- Foreign key to the forms table
    created_at      INT  NOT NULL,               -- User who created the assignment
    access_end_date INT  NOT NULL,               -- End date for access
    CONSTRAINT fk_user FOREIGN KEY (id_user)     -- Foreign key constraint for users
        REFERENCES users (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_role FOREIGN KEY (id_role)     -- Foreign key constraint for roles
        REFERENCES roles (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_module FOREIGN KEY (id_module) -- Foreign key constraint for modules
        REFERENCES modules (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_form FOREIGN KEY (id_form)     -- Foreign key constraint for forms
        REFERENCES forms (id)
        ON DELETE CASCADE
);
CREATE TYPE weekday AS ENUM (
    'Saturday', 'Sunday', 'Monday', 'Tuesday',
    'Wednesday', 'Thursday', 'Friday'
);
CREATE TYPE shift_work AS ENUM (
    'Both', 'Afternoon', 'Morning'
);

CREATE TABLE calendars
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_contract UUID        NOT NULL, -- Foreign key to the contracts table
    shamsi_date VARCHAR(10) NOT NULL, -- Shamsi date as a string (e.g., "1402/05/02")
    work_date   DATE        NOT NULL, -- Gregorian date for the work day
    weekday     weekday     NOT NULL, -- Weekday ENUM type
    year        INT,                  -- Year related to the Shamsi date
    holiday_is  BOOLEAN     NOT NULL, -- Whether the day is a holiday
    shift_work  shift_work  NOT NULL, -- Shift work ENUM type
    description VARCHAR(255),         -- Additional description about the workday
    created_at  INT         NOT NULL,
    CONSTRAINT fk_contract FOREIGN KEY (id_contract)
        REFERENCES contracts (id)
        ON DELETE CASCADE
);

-- Create exceptions table
CREATE TABLE exceptions
(
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each exception
    car_license_plates        TEXT[] DEFAULT ARRAY[]::TEXT[],             -- Array of car license plates included in the exception
    motorcycle_license_plates TEXT[] DEFAULT ARRAY[]::TEXT[],             -- Array of motorcycle license plates included in the exception
    exception_multiplier      DECIMAL(10, 2) NOT NULL,                    -- Multiplier for the exception
    start_date                DATE           NOT NULL,                    -- Start date of the exception validity
    end_date                  DATE,                                       -- End date of the exception validity
    description               TEXT,                                       -- Additional description about the exception
    notification_number       VARCHAR(20),                                -- Notification or letter number
    notification_date         DATE,                                       -- Date of the notification or letter
    document_image            VARCHAR(200),                               -- Uploaded image of the document
    user_id                   UUID           NOT NULL,                    -- Foreign key referencing the users table
    vehicle_type              VARCHAR(20),                                -- Type of vehicle (motorcycle or car)
    created_at                INT            NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (user_id)                              -- Foreign key constraint for user_id
        REFERENCES users (id)
        ON DELETE CASCADE                                                 -- Cascade delete when a user is deleted
);


CREATE TABLE rates
(
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier for each tariff
    code                     VARCHAR(20)    NOT NULL,                    -- Unique code to identify the tariff
    road_category_id         UUID           NOT NULL,                    -- Foreign key referencing road categories
    time_cycle_minutes       INT            NOT NULL,                    -- Duration of each billing cycle in minutes
    rate_multiplier          DECIMAL(10, 2) NOT NULL,                    -- Rate multiplier based on road category and conditions
    peak_hour_multiplier     DECIMAL(10, 2),                             -- Multiplier for peak hours
    good_percentage          INT              DEFAULT 0,                 -- Discount percentage for good customers
    normal_settlement_period INT              DEFAULT 0,                 -- Normal settlement period after receiving debt notification
    late_penalty             DECIMAL(5, 2),                              -- Monthly late penalty percentage
    late_penalty_max         DECIMAL(5, 2),                              -- Maximum late penalty percentage
    valid_from               DATE           NOT NULL,                    -- Start date of tariff validity
    valid_to                 DATE,                                       -- End date of tariff validity
    description              TEXT,                                       -- Additional description
    start_time               VARCHAR(20),                                -- Start time of the tariff
    end_time                 VARCHAR(20),                                -- End time of the tariff
    city_id                  UUID           NOT NULL,                    -- Foreign key referencing cities table
    approval_number          VARCHAR(20),                                -- City council approval number
    approval_date            DATE,                                       -- City council approval date
    year                     INT            NOT NULL,                    -- Tariff year
    base_rate_id             UUID,                                       -- Foreign key referencing base rates table
    exceptions_id            UUID,                                       -- Foreign key referencing exceptions table
    created_at               INT            NOT NULL,

    CONSTRAINT fk_road_category FOREIGN KEY (road_category_id)
        REFERENCES road_categories (id)
        ON DELETE CASCADE,                                               -- Cascade delete when road category is deleted

    CONSTRAINT fk_city FOREIGN KEY (city_id)
        REFERENCES cities (id)
        ON DELETE CASCADE,                                               -- Cascade delete when city is deleted

    CONSTRAINT fk_base_rate FOREIGN KEY (base_rate_id)
        REFERENCES base_rates (id)
        ON DELETE SET NULL,                                              -- Set to NULL if base rate is deleted

    CONSTRAINT fk_exceptions FOREIGN KEY (exceptions_id)
        REFERENCES exceptions (id)
        ON DELETE SET NULL                                               -- Set to NULL if exception is deleted
);