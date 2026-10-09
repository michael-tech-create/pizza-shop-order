-- Pizzas that already appear on orders cannot be hard-deleted (order_items.pizza_id
-- references pizzas with no ON DELETE). is_active lets the admin remove them
-- from the menu while past orders still resolve the pizza name.
ALTER TABLE pizzas ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true;
