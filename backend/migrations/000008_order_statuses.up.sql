CREATE TABLE IF NOT EXISTS order_statuses(
    order_status_id SERIAL PRIMARY KEY,
    order_status_name VARCHAR(100) UNIQUE
);

INSERT INTO order_statuses (order_status_name) VALUES ('Pending'), ('Delivered');