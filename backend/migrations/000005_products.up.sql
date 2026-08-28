CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    product_name VARCHAR(60) NOT NULL,
    product_description VARCHAR(400) NOT NULL,
    product_price DECIMAL(10, 2) NOT NULL,
    product_image VARCHAR(200) NOT NULL,
    product_category_id INTEGER REFERENCES categories(category_id),
    product_seller_id INTEGER REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products(product_seller_id);