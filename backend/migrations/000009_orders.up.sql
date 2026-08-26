CREATE TABLE IF NOT EXISTS orders(
    order_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_status_id INTEGER NOT NULL REFERENCES order_statuses(order_status_id) ON DELETE SET NULL,
    order_total DECIMAL(10, 2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);