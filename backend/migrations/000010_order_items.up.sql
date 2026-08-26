CREATE TABLE IF NOT EXISTS order_items(
    order_item_id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(product_id) ON DELETE SET NULL,
    product_price DECIMAL(10, 2) NOT NULL,
    product_name VARCHAR(60) NOT NULL,
    product_seller_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    order_item_quantity INTEGER NOT NULL,
    order_item_total DECIMAL(10, 2) NOT NULL
);