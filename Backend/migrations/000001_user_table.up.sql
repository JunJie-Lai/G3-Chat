CREATE TABLE IF NOT EXISTS users
(
    id             VARCHAR(50) PRIMARY KEY,
    name          VARCHAR(255)        NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    picture       VARCHAR(255)        NOT NULL,
    refresh_token VARCHAR(255)        NOT NULL
)