-- Insert Payments based on orders, with payment_date = NULL
INSERT INTO payment (user_id, order_id, amount, status)
SELECT 
    o.user_id, 
    o.order_id, 
    o.total_amount,  -- Total amount from orders table
    'pending'        -- Payment status is 'pending'
FROM orders o
WHERE o.order_id IN (1, 2);
