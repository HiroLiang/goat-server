-- Users Table
CREATE TABLE IF NOT EXISTS goat.public.users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    account TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Roles Table
CREATE TABLE IF NOT EXISTS goat.public.roles (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Drop Tables
DROP TABLE IF EXISTS goat.public.users;
DROP TABLE IF EXISTS goat.public.roles;
