CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE users (
                       user_id            UUID   DEFAULT uuid_generate_v4() ,
                       msisdn VARCHAR(15),
                       msisdn_verified BOOLEAN NOT NULL DEFAULT FALSE,
                       encrypted_password VARCHAR(256) NOT NULL,
                       contractor_name VARCHAR(100) NOT NULL,
                       contractor_code VARCHAR(10),
                       registration_number VARCHAR(50),
                       contact_person VARCHAR(100),
                       ceo_name VARCHAR(100),
                       authorized_signatories TEXT,
                       phone_number VARCHAR(15),
                       email VARCHAR(100),
                       address VARCHAR(255),
                       contract_type VARCHAR(50),
                       bank_account_number VARCHAR(30),
                       role VARCHAR(30),
                       description TEXT,
                       created_at INT,
                       updated_at INT
);



CREATE TABLE contractors
(
    id                     SERIAL PRIMARY KEY,
    contractor_name        VARCHAR(100)       NOT NULL,
    code_contractor        VARCHAR(10) UNIQUE NOT NULL,
    number_registration    VARCHAR(50),
    person_contact         VARCHAR(100),
    ceo_name               VARCHAR(100),
    signatories_authorized TEXT,
    phone_number           VARCHAR(15),
    email                  VARCHAR(100),
    address                VARCHAR(255),
    type_contract          VARCHAR(50),
    number_account_bank    VARCHAR(30),
    description            TEXT
);


CREATE TABLE license_plate_reader_devices
(
    id                    SERIAL PRIMARY KEY,
    code_device           VARCHAR(20) NOT NULL,
    number_serial         VARCHAR(50) NOT NULL,
    model                 VARCHAR(50),
    date_installation     DATE,
    date_expiry_warranty  DATE,
    date_expiry_insurance DATE,
    class_device          VARCHAR(50),
    image_contract        BYTEA,
    image_insurance       BYTEA,
    id_contractor         INT REFERENCES contractors (id),
    description           TEXT
);

CREATE TABLE vehicles
(
    id                           SERIAL PRIMARY KEY,              -- Unique identifier for each vehicle
    code_vehicle                 VARCHAR(20) UNIQUE NOT NULL,     -- Unique code for identifying the vehicle
    vin                          VARCHAR(20) UNIQUE NOT NULL,     -- Vehicle Identification Number (VIN)
    plate_license                VARCHAR(15) UNIQUE NOT NULL,     -- License plate number
    type_vehicle                 VARCHAR(50),                     -- Type of vehicle (e.g., motorcycle, car, minibus, truck, etc.)
    brand                        VARCHAR(50),                     -- Brand of the vehicle (e.g., Peugeot, Hyundai)
    model                        VARCHAR(50),                     -- Model of the vehicle (e.g., 206, Sonata)
    color                        VARCHAR(30),                     -- Color of the vehicle
    manufacture_of_year          INT,                             -- Manufacturing year of the vehicle
    kilometers_initial           BIGINT,                          -- Initial kilometers at the start of activity
    expiry_insurance_party_third DATE,                            -- Expiry date of third-party insurance
    expiry_insurance_body        DATE,                            -- Expiry date of body insurance
    image_document_vehicle       BYTEA,                           -- Image of the vehicle's document
    image_card_vehicle           BYTEA,                           -- Image of the vehicle's card
    third_party_insurance_image  BYTEA,                           -- Image of third-party insurance
    body_insurance_image         BYTEA,                           -- Image of body insurance
    id_contractor                INT REFERENCES contractors (id), -- Foreign key to the contractors table
    status                       VARCHAR(20),                     -- Status of the vehicle (e.g., active, inactive, in repair)
    description                  TEXT                             -- Additional description or notes
);

CREATE TABLE lpr_vehicles_combination
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lpr_id     INT REFERENCES license_plate_reader_devices (id),
    vehicle_ID INT REFERENCES vehicles (id)
)