CREATE SCHEMA IF NOT EXISTS authserver;

CREATE TABLE IF NOT EXISTS authserver.user_profiles
(
    username      TEXT PRIMARY KEY,
    binds         TEXT[],
    image         TEXT,
    networks      TEXT[],
    public_keys   TEXT[],
    shell_command TEXT[],
    secret        TEXT
);
