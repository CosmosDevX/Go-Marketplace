CREATE TABLE IF NOT EXISTS roles(
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO roles (role_name) VALUES ('user'), ('seller'), ('admin');