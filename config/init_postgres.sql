-- User Table
CREATE TABLE IF NOT EXISTS goat.public.users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    account TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Role
CREATE TABLE IF NOT EXISTS goat.public.roles (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    permissions TEXT NOT NULL,
    name TEXT NOT NULL
);

-- Features
CREATE TABLE IF NOT EXISTS goat.public.features (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Users-Roles
CREATE TABLE IF NOT EXISTS goat.public.user_role (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

--Role-Features
CREATE TABLE IF NOT EXISTS goat.public.role_feature (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    feature_id BIGINT NOT NULL REFERENCES features(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, feature_id)
);

-- Chat room
CREATE TABLE IF NOT EXISTS goat.public.chat_rooms (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Chat records (role to role)
CREATE TABLE IF NOT EXISTS goat.public.chat_record (
    from_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chat_room_id BIGINT NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (from_id, to_id, chat_room_id)
);

-- Refresh token
CREATE TABLE IF NOT EXISTS goat.public.refresh_tokens (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    expired_at TIMESTAMP NOT NULL DEFAULT now() + interval '1 day',
    PRIMARY KEY (user_id, token)
);

-- llama
-- CREATE TABLE IF NOT EXISTS goat.public.llama ();

-- Role-Llamas (for prompt. maybe in the other table)
-- CREATE TABLE IF NOT EXISTS goat.public.role_llama ();

-- Session of llama
-- CREATE TABLE IF NOT EXISTS goat.public.chat_session ();

-- Drop Tables
DROP TABLE IF EXISTS goat.public.role_llama;
DROP TABLE IF EXISTS goat.public.chat_record;
DROP TABLE IF EXISTS goat.public.chat_rooms;
DROP TABLE IF EXISTS goat.public.role_feature;
DROP TABLE IF EXISTS goat.public.user_role;
DROP TABLE IF EXISTS goat.public.chat_session;
DROP TABLE IF EXISTS goat.public.refresh_tokens;

DROP TABLE IF EXISTS goat.public.llamas;
DROP TABLE IF EXISTS goat.public.features;
DROP TABLE IF EXISTS goat.public.roles;
DROP TABLE IF EXISTS goat.public.users;
