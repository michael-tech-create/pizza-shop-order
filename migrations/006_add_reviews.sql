-- One review per order — the UNIQUE constraint on order_id is what
-- prevents someone from submitting multiple reviews against the same
-- verified order.
CREATE TABLE IF NOT EXISTS reviews (
    id          SERIAL PRIMARY KEY,
    order_id    INT NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    customer_id INT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    rating      INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reviews_customer_id ON reviews(customer_id);
CREATE INDEX IF NOT EXISTS idx_reviews_order_id ON reviews(order_id);