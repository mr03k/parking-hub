CREATE TABLE roles
(
    id         UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    title      VARCHAR(100) NOT NULL UNIQUE,
    created_at BIGINT       NOT NULL,
    updated_at BIGINT       NOT NULL,
    deleted_at BIGINT
);


-- Insert admin role
INSERT INTO roles (title, created_at, updated_at,deleted_at)
VALUES ('Admin', extract(epoch from now()), extract(epoch from now()),0);

-- Insert driver role
INSERT INTO roles (title, created_at, updated_at,deleted_at)
VALUES ('Driver', extract(epoch from now()), extract(epoch from now()),0);