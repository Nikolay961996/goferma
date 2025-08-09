CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    order_number TEXT NOT NULL,
    accrual BIGINT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT unique_user_order UNIQUE (user_id, order_number)
);
