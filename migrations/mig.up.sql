DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS product_categories CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS cart_items CASCADE;

CREATE TABLE IF NOT EXISTS product_categories(
	id SERIAL PRIMARY KEY,
	name VARCHAR UNIQUE NOT NULL,
	description TEXT
);

CREATE TABLE IF NOT EXISTS products(
	id SERIAL PRIMARY KEY,
	name VARCHAR NOT NULL,
	description TEXT,
	price DECIMAL(10,2) NOT NULL,
	quantity INTEGER NOT NULL,
	category_id INTEGER REFERENCES product_categories(id),
	image_path VARCHAR
);

CREATE TABLE IF NOT EXISTS users(
	id SERIAL PRIMARY KEY,
	first_name VARCHAR NOT NULL,
	last_name VARCHAR NOT NULL,
	email VARCHAR UNIQUE NOT NULL,
	password VARCHAR NOT NULL,
	phone VARCHAR UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS orders(
	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
	date DATE NOT NULL,
	price DECIMAL(10,2) NOT NULL,
	status VARCHAR(20) CHECK (status IN ('pending', 'paid', 'canceled')) NOT NULL,
	pickup_method VARCHAR(20) CHECK (pickup_method IN ('delivery', 'pickup')) NOT NULL,
	delivery_address VARCHAR
);

CREATE TABLE IF NOT EXISTS order_items(
		order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
	product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
	quantity INTEGER NOT NULL,
	product_price DECIMAL(10,2) NOT NULL,
	CONSTRAINT order_items_pk PRIMARY KEY(order_id, product_id)
);

CREATE TABLE IF NOT EXISTS payments(
	id SERIAL PRIMARY KEY,
	order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
	payment DECIMAL(10,2) NOT NULL,
	date DATE NOT NULL
); 

CREATE TABLE IF NOT EXISTS cart_items(
	user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
	quantity INTEGER NOT NULL,
	CONSTRAINT carts_pk PRIMARY KEY(user_id, product_id)
);