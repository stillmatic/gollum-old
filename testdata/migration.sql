-- Creating the 'customers' table
CREATE TABLE customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone TEXT
);

-- Inserting test data into the 'customers' table
INSERT INTO customers (first_name, last_name, email, phone) VALUES
('John', 'Doe', 'john.doe@example.com', '555-123-4567'),
('Jane', 'Smith', 'jane.smith@example.com', '555-987-6543'),
('Michael', 'Johnson', 'michael.johnson@example.com', '555-234-5678'),
('Emily', 'Williams', 'emily.williams@example.com', '555-876-5432');

-- Creating the 'orders' table
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_id INTEGER NOT NULL,
    order_date TEXT NOT NULL,
    total_amount REAL NOT NULL,
    status TEXT NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers (id)
);

-- Inserting test data into the 'orders' table
INSERT INTO orders (customer_id, order_date, total_amount, status) VALUES
(1, '2023-01-12', 150.50, 'shipped'),
(1, '2023-02-14', 87.90, 'delivered'),
(2, '2023-01-25', 230.00, 'shipped'),
(2, '2023-02-28', 130.50, 'delivered'),
(3, '2023-01-30', 100.00, 'shipped'),
(4, '2023-02-02', 75.00, 'processing');