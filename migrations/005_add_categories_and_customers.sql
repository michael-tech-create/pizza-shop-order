-- ============================================================
-- Categories (fully admin-manageable)
-- ============================================================
CREATE TABLE IF NOT EXISTS categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Nullable — existing pizzas keep working uncategorized until an admin
-- assigns them one.
ALTER TABLE pizzas ADD COLUMN IF NOT EXISTS category_id INT REFERENCES categories(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_pizzas_category_id ON pizzas(category_id);

-- ============================================================
-- Customer accounts (optional — guest checkout still fully works)
-- Login is phone + password, matching the rest of this app's
-- phone-first design instead of requiring an email.
-- ============================================================
CREATE TABLE IF NOT EXISTS customers (
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(150) NOT NULL,
    phone         VARCHAR(20) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Nullable — a guest order simply has customer_id = NULL, exactly like
-- every order placed before this migration.
ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_id INT REFERENCES customers(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);

-- Starter categories — edit/add more from the admin dashboard afterward.
INSERT INTO categories (name) VALUES
    ('Classic'), ('Specialty'), ('Vegetarian'), ('Meat Lovers')
ON CONFLICT DO NOTHING;