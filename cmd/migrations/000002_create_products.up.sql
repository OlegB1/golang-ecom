CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image TEXT NOT NULL,
    price INTEGER NOT NULL, 
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);