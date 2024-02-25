CREATE TABLE IF NOT EXISTS users
(
    id           bigserial PRIMARY KEY,
    createdAt    timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt    timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    profilePhoto text                        NOT NULL,
    name         text                        NOT NULL,
    username     text                        NOT NULL,
    description  text                        NOT NULL,
    email        text                        NOT NULL,
    password     text                        NOT NULL
);