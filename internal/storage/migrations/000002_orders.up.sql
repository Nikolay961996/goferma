CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    order_number TEXT NOT NULL,
    CONSTRAINT unique_user_order UNIQUE (user_id, order_number)
);
