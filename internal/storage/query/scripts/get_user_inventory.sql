SELECT m.name,i.quantity
FROM inventory i
JOIN merch m ON i.merch_id = m.id
WHERE i.user_id = $1;