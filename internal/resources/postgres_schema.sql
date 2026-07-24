CREATE SCHEMA IF NOT EXISTS authserver;

CREATE TABLE IF NOT EXISTS authserver.user_profiles
(
    username      VARCHAR(32) PRIMARY KEY,
    binds         TEXT[],
    image         VARCHAR(384),
    networks      VARCHAR(64)[],
    public_keys   TEXT[],
    shell_command TEXT[],
    secret        VARCHAR(255)
);
