INSERT INTO inventory (user_id,merch_id,quantity)
VALUES ($1,$2,1)
ON CONFLICT (user_id,merch_id)
DO UPDATE SET quantity = inventory.quantity +1;