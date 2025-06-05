-- Insert default admin staff
INSERT INTO staff (full_name, position, email, password_hash, created_at)
VALUES ('Admin', 'admin', 'admin@uniclever.com', 'admin123', CURRENT_TIMESTAMP);
