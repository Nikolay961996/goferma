CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL,
    passwordHash TEXT NOT NULL
);
