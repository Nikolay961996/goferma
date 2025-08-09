CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL,
    password_hash TEXT NOT NULL
);
