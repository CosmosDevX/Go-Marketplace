CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    product_name VARCHAR(60) NOT NULL,
    product_description VARCHAR(400) NOT NULL,
    product_price DECIMAL(10, 2) NOT NULL,
    product_quantity INTEGER DEFAULT 0,
    product_image VARCHAR(200) NOT NULL,
    product_category_id INT REFERENCES categories(category_id)
);