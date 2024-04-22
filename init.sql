CREATE TYPE allowance_type AS ENUM ('personal', 'donation', 'k-receipt');

CREATE TABLE IF NOT EXISTS allowances (
  id SERIAL PRIMARY KEY,
  type allowance_type NOT NULL,
  init_amount DECIMAL(10, 2) NOT NULL,
  min_amount DECIMAL(10, 2) NOT NULL,
  max_amount DECIMAL(10, 2) NOT NULL,
  limit_max_amount DECIMAL(10, 2) NOT NULL, 
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


INSERT INTO allowances (type, init_amount,min_amount, max_amount, limit_max_amount) VALUES 
('personal', 60000, 10000.00, 100000.00, 100000.00), 
('donation', 0, 0, 100000.00, 100000.00), 
('k-receipt', 0, 0, 50000.00, 100000.00);