CREATE TABLE pos_cash_drawer (
    drawer_id UUID PRIMARY KEY,
    store_id UUID,
    employee_id UUID,
    cash_in DECIMAL(10, 2) NOT NULL,
    cash_out DECIMAL(10, 2) NOT NULL,
    transaction_time TIMESTAMP NOT NULL,
    role_id UUID REFERENCES pos_roles(role_id),
    branch_id UUID NOT NULL,
    company_id UUID NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_customers (
    customer_id UUID PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone_number VARCHAR(20),
    date_of_birth DATE,
    registration_date DATE,
    address VARCHAR(255),
    city VARCHAR(100),
    country VARCHAR(100),
    branch_id UUID NOT NULL,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_invoices (
    invoice_id UUID PRIMARY KEY,
    sale_id UUID REFERENCES pos_sales(sale_id),
    date TIMESTAMP NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    discounts DECIMAL(10, 2),
    taxes DECIMAL(10, 2),
    branch_id UUID NOT NULL,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_online_payments (
    payment_id UUID PRIMARY KEY,
    store_id UUID,
    sale_id UUID,
    employee_id UUID,
    payment_date TIMESTAMP NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_method UUID,
    role_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_payment_methods (
    payment_method_id UUID PRIMARY KEY,
    method_name VARCHAR(255) NOT NULL,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_returns (
    return_id UUID PRIMARY KEY,
    sale_id UUID REFERENCES pos_sales(sale_id),
    product_id UUID NOT NULL,
    quantity INT NOT NULL,
    return_date TIMESTAMP NOT NULL,
    reason TEXT,
    branch_id UUID NOT NULL,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);

CREATE TABLE pos_sales (
    sale_id UUID PRIMARY KEY,
    receipt_id VARCHAR(255) NOT NULL,
    product_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    sale_date TIMESTAMP NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    store_id UUID,
    cashier_id UUID,
    payment_method_id UUID NOT NULL,
    branch_id UUID NOT NULL,
    company_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP,
    updated_by UUID
);



-- Memasukkan data ke dalam pos_customers
INSERT INTO pos_customers (customer_id, first_name, last_name, email, phone_number, date_of_birth, registration_date, address, city, country, company_id) VALUES
('111e8400-e29b-41d4-a716-446655440000', 'Budi', 'Santoso', 'budi.santoso@gmail.com', '081234567890', '1980-01-01', '2024-05-01', 'Jl. Jend. Sudirman No.1', 'Jakarta', 'Indonesia', '550e8400-e29b-41d4-a716-446655440000'),
('111e8400-e29b-41d4-a716-446655440001', 'Susi', 'Widianti', 'susi.widianti@gmail.com', '081234567891', '1985-01-01', '2024-05-01', 'Jl. MH. Thamrin No.1', 'Jakarta', 'Indonesia', '550e8400-e29b-41d4-a716-446655440000');

-- Memasukkan data ke dalam pos_cash_drawer
INSERT INTO pos_cash_drawer (drawer_id, store_id, employee_id, cash_in, cash_out, transaction_time, role_id, company_id) VALUES
('222e8400-e29b-41d4-a716-446655440000', '330e8400-e29b-41d4-a716-446655440000', '444e8400-e29b-41d4-a716-446655440000', 1000000.00, 500000.00, '2024-05-01 09:00:00', '555e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000');

-- Memasukkan data ke dalam pos_online_payments
INSERT INTO pos_online_payments (payment_id, sale_id, employee_id, payment_date, amount, payment_method, role_id, company_id) VALUES
('333e8400-e29b-41d4-a716-446655440000', '777e8400-e29b-41d4-a716-446655440000', '444e8400-e29b-41d4-a716-446655440000', '2024-05-01 09:30:00', 500000.00, '666e8400-e29b-41d4-a716-446655440000', '555e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000');

-- Memasukkan data ke dalam pos_sales
INSERT INTO pos_sales (sale_id, product_id, customer_id, quantity, sale_date, total_price, store_id, cashier_id, payment_method_id, company_id) VALUES
('777e8400-e29b-41d4-a716-446655440000', '770e8400-e29b-41d4-a716-446655440000', '111e8400-e29b-41d4-a716-446655440000', 5, '2024-05-01 09:30:00', 100000.00, '330e8400-e29b-41d4-a716-446655440000', '444e8400-e29b-41d4-a716-446655440000', '666e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000');

-- Memasukkan data ke dalam pos_payment_methods
INSERT INTO pos_payment_methods (payment_method_id, method_name, company_id) VALUES
('666e8400-e29b-41d4-a716-446655440000', 'Kartu Kredit', '550e8400-e29b-41d4-a716-446655440000'),
('666e8400-e29b-41d4-a716-446655440001', 'Transfer Bank', '550e8400-e29b-41d4-a716-446655440000');

-- Memasukkan data ke dalam pos_invoices
INSERT INTO pos_invoices (invoice_id, sale_id, date, total, discounts, taxes, company_id) VALUES
('888e8400-e29b-41d4-a716-446655440000', '777e8400-e29b-41d4-a716-446655440000', '2024-05-01 09:30:00', 100000.00, 0.00, 10000.00, '550e8400-e29b-41d4-a716-446655440000');

-- Memasukkan data ke dalam pos_returns
INSERT INTO pos_returns (return_id, sale_id, product_id, quantity, return_date, reason, company_id) VALUES
('999e8400-e29b-41d4-a716-446655440000', '777e8400-e29b-41d4-a716-446655440000', '770e8400-e29b-41d4-a716-446655440000', 1, '2024-05-02 09:30:00', 'Produk rusak', '550e8400-e29b-41d4-a716-446655440000');
