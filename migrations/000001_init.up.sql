CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(32) PRIMARY KEY,
    data jsonb NOT NULL
);
