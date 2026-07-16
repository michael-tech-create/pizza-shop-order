-- ============================================================
-- McMichael Pizza — full schema for a fresh Supabase database
-- Run this once against your new (empty) Supabase Postgres.
-- Reconstructed from everything the Go backend expects — if your
-- original local DB had anything extra not reflected in code
-- (e.g. additional constraints), add those separately afterward.
-- ============================================================

-- ADMINS
CREATE TABLE IF NOT EXISTS admins (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- PIZZAS (per-size pricing — small/medium/large)
CREATE TABLE IF NOT EXISTS pizzas (
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(150) NOT NULL,
    price_small  INT NOT NULL,
    price_medium INT NOT NULL,
    price_large  INT NOT NULL,
    description  TEXT
);

-- PIZZA IMAGES
CREATE TABLE IF NOT EXISTS pizza_images (
    id        SERIAL PRIMARY KEY,
    pizza_id  INT NOT NULL REFERENCES pizzas(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pizza_images_pizza_id ON pizza_images(pizza_id);

-- ORDERS (delivery details + payment tracking)
CREATE TABLE IF NOT EXISTS orders (
    id                 SERIAL PRIMARY KEY,
    customer_name      VARCHAR(150) NOT NULL,
    phone              VARCHAR(20) NOT NULL,
    address            TEXT NOT NULL,
    total_cost         INT NOT NULL DEFAULT 0,
    status             VARCHAR(20) NOT NULL DEFAULT 'pending',       -- pending | preparing | delivered | cancelled
    payment_status     VARCHAR(20) NOT NULL DEFAULT 'pending',       -- pending | paid | failed
    payment_reference  VARCHAR(100)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_payment_reference
    ON orders(payment_reference)
    WHERE payment_reference IS NOT NULL;

-- ORDER ITEMS (one row per pizza+size+quantity within an order)
CREATE TABLE IF NOT EXISTS order_items (
    id        SERIAL PRIMARY KEY,
    order_id  INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    pizza_id  INT NOT NULL REFERENCES pizzas(id),
    size      VARCHAR(10) NOT NULL DEFAULT 'medium',                 -- small | medium | large
    quantity  INT NOT NULL,
    sub_total INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_pizza_id ON order_items(pizza_id);

-- ============================================================
-- Optional: seed a couple of pizzas so /menu isn't empty on first load.
-- Comment this block out if you'd rather add pizzas via the admin
-- dashboard instead.
-- ============================================================
INSERT INTO pizzas (name, price_small, price_medium, price_large, description)
VALUES
    ('Pepperoni',  3600, 4500, 5850, 'Spicy pepperoni with extra cheese'),
    ('BBQ Chicken', 3600, 4500, 5850, 'Original smoky BBQ flavor')
ON CONFLICT DO NOTHING;
