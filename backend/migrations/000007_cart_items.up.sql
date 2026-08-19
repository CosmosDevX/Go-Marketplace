CREATE TABLE IF NOT EXISTS cart_items(
    cart_item_id SERIAL PRIMARY KEY,
    cart_id INTEGER NOT NULL REFERENCES carts(cart_id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
    UNIQUE(cart_id, product_id),
    quantity INTEGER DEFAULT 1
);