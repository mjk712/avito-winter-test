SELECT
    th.transaction_type,
    th.amount,
    th.timestamp,
    u1.username AS from_user,
    u2.username AS to_user,
    m.name AS merch_name
FROM transaction_history th
LEFT JOIN users u1 ON th.from_user_id = u1.id
JOIN users u2 ON th.to_user_id = u2.id
LEFT JOIN merch m ON th.merch_id = m.id
WHERE th.to_user_id = $1 OR th.from_user_id = $1
ORDER BY th.timestamp DESC;