-- Northwind-like schema for NeuronSQL eval
CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  category_id INT,
  unit_price DECIMAL(10,2)
);
CREATE TABLE IF NOT EXISTS categories (id SERIAL PRIMARY KEY, name TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS orders (id SERIAL PRIMARY KEY, product_id INT, quantity INT, created_at TIMESTAMPTZ DEFAULT NOW());
