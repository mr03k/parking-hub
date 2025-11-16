CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
EXTENSION IF NOT EXISTS postgis;

CREATE TABLE users
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL,
    password      VARCHAR(255) NOT NULL,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    email         VARCHAR(100) NOT NULL UNIQUE,
    phone_number  VARCHAR(15)  NOT NULL UNIQUE,
    national_id   VARCHAR(50),
    postal_code   VARCHAR(50),
    company_name  VARCHAR(100),
    profile_image VARCHAR(150),
    gender        VARCHAR(6),
    address       VARCHAR(255),
    status        VARCHAR(8),
    role          VARCHAR(14),
    created_at    BIGINT       NOT NULL,
    updated_at    BIGINT       NOT NULL,
    deleted_at    BIGINT
);

CREATE TABLE contractors
(
    id                     UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
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
    description            TEXT,
    created_at             BIGINT             NOT NULL,
    updated_at             BIGINT             NOT NULL,
    deleted_at             BIGINT
);

CREATE TABLE contracts
(
    id               UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    contract_number  VARCHAR(50) NOT NULL UNIQUE,
    contract_date    DATE        NOT NULL,
    start_date       DATE        NOT NULL,
    end_date         DATE        NOT NULL,
    contract_amount  BIGINT      NOT NULL,
    contract_type    VARCHAR(50) NOT NULL,
    contractor_id    UUID        NOT NULL REFERENCES contractors (id) ON DELETE CASCADE,
    operation_period INT         NOT NULL,
    equipment_period INT         NOT NULL,
    description      TEXT,
    created_at       BIGINT      NOT NULL,
    updated_at       BIGINT      NOT NULL,
    deleted_at       BIGINT
);


CREATE TABLE vehicles
(
    id                           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    vehicle_id                   BIGINT UNIQUE       NOT NULL,
    code_vehicle                 VARCHAR(120) UNIQUE NOT NULL,
    vin                          VARCHAR(120) UNIQUE NOT NULL,
    plate_license                VARCHAR(55) UNIQUE  NOT NULL,
    type_vehicle                 VARCHAR(50),
    brand                        VARCHAR(50),
    model                        VARCHAR(50),
    color                        VARCHAR(30),
    manufacture_of_year          INT,
    kilometers_initial           BIGINT,
    expiry_insurance_party_third BIGINT,
    expiry_insurance_body        BIGINT,
    image_document_vehicle       VARCHAR(256),
    image_card_vehicle           VARCHAR(256),
    third_party_insurance_image  VARCHAR(256),
    body_insurance_image         VARCHAR(256),
    contractor_id                UUID                NOT NULL REFERENCES contractors (id) ON DELETE CASCADE,
    status                       VARCHAR(20),
    description                  TEXT,
    created_at                   BIGINT              NOT NULL,
    updated_at                   BIGINT              NOT NULL,
    deleted_at                   BIGINT
);

CREATE TABLE devices
(
    id                    UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    device_id             BIGINT UNIQUE NOT NULL,
    code_device           VARCHAR(120)  NOT NULL,
    number_serial         VARCHAR(50)   NOT NULL,
    model                 VARCHAR(50),
    date_installation     BIGINT,
    date_expiry_warranty  BIGINT,
    date_expiry_insurance BIGINT,
    class_device          VARCHAR(50),
    image_contract        VARCHAR(256),
    image_insurance       VARCHAR(256),
    contractor_id         UUID          NOT NULL REFERENCES contractors (id) ON DELETE CASCADE,
    vehicle_id            UUID          NOT NULL REFERENCES vehicles (id) ON DELETE CASCADE,
    description           TEXT,
    created_at            BIGINT        NOT NULL,
    updated_at            BIGINT        NOT NULL,
    deleted_at            BIGINT
);

CREATE TABLE drivers
(
    id                          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    contractor_id               UUID   REFERENCES contractors (id) ON DELETE SET NULL,
    user_id                     UUID REFERENCES users (id) ON DELETE CASCADE,
    driver_type                 VARCHAR(10),
    shift_type                  VARCHAR(10),
    employment_status           VARCHAR(20),
    employment_start_date       BIGINT,
    employment_end_date         BIGINT,
    driver_photo                VARCHAR(256),
    id_card_image               VARCHAR(256),
    birth_certificate_image     VARCHAR(256),
    military_service_card_image VARCHAR(256),
    health_certificate_image    VARCHAR(256),
    criminal_record_image       VARCHAR(256),
    description                 VARCHAR(256),
    created_at                  BIGINT NOT NULL,
    updated_at                  BIGINT NOT NULL,
    deleted_at                  BIGINT
);



CREATE TABLE calenders
(
    id               UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    contract_id      UUID REFERENCES contracts (id) ON DELETE CASCADE,
    shamsi_date      VARCHAR(10) NOT NULL,
    work_date        BIGINT      NOT NULL,
    weekday          VARCHAR,
    year             INT         NOT NULL,
    is_holiday       BOOLEAN     NOT NULL,
    work_shift       VARCHAR(17),
    description      TEXT,
    work_shift_start BIGINT,
    work_shift_end   BIGINT,
    created_at       BIGINT      NOT NULL,
    updated_at       BIGINT      NOT NULL,
    deleted_at       BIGINT

);
CREATE TABLE rings
(
    id         BIGINT PRIMARY KEY,
    ring_code  VARCHAR(100),
    length     double precision,
    ring_name  VARCHAR(120),
    geom       GEOMETRY,

    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT

);



CREATE TABLE driver_assignments
(
    id           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    driver_id    UUID         NOT NULL,
    code_vehicle VARCHAR(150) NOT NULL UNIQUE,
    ring_id      BIGINT       NOT NULL,
    calender_id  UUID         NOT NULL,
    created_at   BIGINT       NOT NULL,
    updated_at   BIGINT       NOT NULL,
    deleted_at   BIGINT,

    CONSTRAINT fk_driver FOREIGN KEY (driver_id) REFERENCES drivers (id) ON DELETE CASCADE,
    CONSTRAINT fk_vehicle FOREIGN KEY (code_vehicle) REFERENCES vehicles (code_vehicle) ON DELETE CASCADE,
    CONSTRAINT fk_ring FOREIGN KEY (ring_id) REFERENCES rings (id) ON DELETE CASCADE,
    CONSTRAINT fk_calender FOREIGN KEY (calender_id) REFERENCES calenders (id) ON DELETE CASCADE
);


CREATE TABLE roads
(
    id          BIGINT PRIMARY KEY,
    geom        GEOMETRY,
    road_name   VARCHAR(100)     NOT NULL,
    road_code   INTEGER          NOT NULL,
    description VARCHAR(200),

    length      DOUBLE PRECISION,
    speed_limit INTEGER, -- Crop of plate photo

    road_type   VARCHAR(100)     NOT NULL,
    road_grade  VARCHAR(1)       NOT NULL,

    created_at  BIGINT           NOT NULL,
    updated_at  BIGINT           NOT NULL,
    deleted_at  BIGINT
);


CREATE TABLE segments
(
    id          BIGINT PRIMARY KEY,
    geom        GEOMETRY,
    seg_name    VARCHAR(200),
    seg_code    VARCHAR(100),
    description VARCHAR(200),
    seg_length  double precision,

    created_at  BIGINT NOT NULL,
    updated_at  BIGINT NOT NULL,
    deleted_at  BIGINT
);

CREATE TABLE parkings
(
    id          BIGINT PRIMARY KEY,
    geom        GEOMETRY,
    park_code   integer,
    park_type   VARCHAR(1) NOT NULL,
    position    VARCHAR(1) NOT NULL,
    description VARCHAR(200),

    created_at  BIGINT     NOT NULL,
    updated_at  BIGINT     NOT NULL,
    deleted_at  BIGINT
);

