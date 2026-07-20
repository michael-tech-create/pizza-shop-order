-- Adds payment tracking to orders.
-- payment_status is separate from `status` (which tracks kitchen/delivery
-- progress) because a paid order can still be "pending" in the kitchen,
-- and an unpaid order should never move to "preparing".
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_status VARCHAR(20) NOT NULL DEFAULT 'pending';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_reference VARCHAR(100);

-- Speeds up the webhook/verify lookups, which search by reference.
CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_payment_reference
    ON orders(payment_reference)
    WHERE payment_reference IS NOT NULL;