-- DBMS: PostgreSQL

DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS fare_matrix CASCADE;
DROP TABLE IF EXISTS cards CASCADE;
DROP TABLE IF EXISTS gates CASCADE;
DROP TABLE IF EXISTS terminals CASCADE;

DROP TYPE IF EXISTS transaction_type CASCADE;
DROP TYPE IF EXISTS card_status CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE terminals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE gates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(10) NOT NULL,
    name VARCHAR(50) NOT NULL,
    terminal_id UUID NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_gates_terminal FOREIGN KEY (terminal_id) REFERENCES terminals(id) ON DELETE RESTRICT,
    CONSTRAINT unique_gate_code_per_terminal UNIQUE (terminal_id, code)
);

CREATE TABLE fare_matrix (
    origin_terminal_id UUID NOT NULL,
    destination_terminal_id UUID NOT NULL,
    fare_amount NUMERIC(8, 2) NOT NULL CHECK (fare_amount >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (origin_terminal_id, destination_terminal_id),
    CONSTRAINT fk_fare_origin FOREIGN KEY (origin_terminal_id) REFERENCES terminals(id) ON DELETE CASCADE,
    CONSTRAINT fk_fare_destination FOREIGN KEY (destination_terminal_id) REFERENCES terminals(id) ON DELETE CASCADE
);

CREATE TYPE card_status AS ENUM ('active', 'blocked', 'expired');

CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    card_number VARCHAR(16) NOT NULL UNIQUE,
    balance NUMERIC(8, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    status card_status NOT NULL DEFAULT 'active',
    issued_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expiry_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE transaction_type AS ENUM ('tap_in', 'tap_out');

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    card_id UUID NOT NULL,
    gate_id UUID NOT NULL,
    terminal_id UUID NOT NULL,
    transaction_type transaction_type NOT NULL,
    amount NUMERIC(8, 2),
    transaction_time TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_trans_card FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE RESTRICT,
    CONSTRAINT fk_trans_gate FOREIGN KEY (gate_id) REFERENCES gates(id) ON DELETE RESTRICT,
    CONSTRAINT fk_trans_terminal FOREIGN KEY (terminal_id) REFERENCES terminals(id) ON DELETE RESTRICT
);

CREATE INDEX idx_gates_terminal ON gates(terminal_id);
CREATE INDEX idx_transactions_card ON transactions(card_id);
CREATE INDEX idx_transactions_time ON transactions(transaction_time);
CREATE INDEX idx_cards_number ON cards(card_number);

CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE INDEX idx_admin_username ON admins(username);

INSERT INTO terminals (id, code, name, address, is_active) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'TRM001', 'Central Terminal', '123 Main Street, Downtown Business District', true),
('550e8400-e29b-41d4-a716-446655440001', 'TRM002', 'Airport Terminal', '456 Airport Road, International Airport', true),
('550e8400-e29b-41d4-a716-446655440002', 'TRM003', 'North Station', '789 North Avenue, North City Center', true),
('550e8400-e29b-41d4-a716-446655440003', 'TRM004', 'South Hub', '321 South Boulevard, South District', true),
('550e8400-e29b-41d4-a716-446655440004', 'TRM005', 'East Terminal', '654 East Street, East Commercial Zone', true),
('550e8400-e29b-41d4-a716-446655440005', 'TRM006', 'West Station', '987 West Road, West Residential Area', true),
('550e8400-e29b-41d4-a716-446655440006', 'TRM007', 'University Terminal', '147 Campus Drive, University District', true),
('550e8400-e29b-41d4-a716-446655440007', 'TRM008', 'Shopping Mall Hub', '258 Mall Avenue, Shopping Center', true),
('550e8400-e29b-41d4-a716-446655440008', 'TRM009', 'Industrial Zone', '369 Factory Road, Industrial District', true),
('550e8400-e29b-41d4-a716-446655440009', 'TRM010', 'Hospital Station', '741 Medical Center Drive, Healthcare District', true),
('550e8400-e29b-41d4-a716-44665544000a', 'TRM011', 'Sports Complex', '852 Stadium Street, Sports District', true),
('550e8400-e29b-41d4-a716-44665544000b', 'TRM012', 'Beach Terminal', '963 Coastal Highway, Beach Area', true),
('550e8400-e29b-41d4-a716-44665544000c', 'TRM013', 'Mountain Station', '159 Highland Road, Mountain Resort', true),
('550e8400-e29b-41d4-a716-44665544000d', 'TRM014', 'Tech Park Hub', '357 Innovation Boulevard, Technology Park', true),
('550e8400-e29b-41d4-a716-44665544000e', 'TRM015', 'Old Town Terminal', '468 Heritage Street, Historic District', false);

INSERT INTO gates (id, code, name, terminal_id, is_active) VALUES
-- Central Terminal Gates
('660e8400-e29b-41d4-a716-446655440000', 'GT001', 'Entry Gate A', '550e8400-e29b-41d4-a716-446655440000', true),
('660e8400-e29b-41d4-a716-446655440001', 'GT002', 'Entry Gate B', '550e8400-e29b-41d4-a716-446655440000', true),
('660e8400-e29b-41d4-a716-446655440002', 'GT003', 'Exit Gate A', '550e8400-e29b-41d4-a716-446655440000', true),
-- Airport Terminal Gates
('660e8400-e29b-41d4-a716-446655440003', 'GT004', 'Entry Gate 1', '550e8400-e29b-41d4-a716-446655440001', true),
('660e8400-e29b-41d4-a716-446655440004', 'GT005', 'Entry Gate 2', '550e8400-e29b-41d4-a716-446655440001', true),
('660e8400-e29b-41d4-a716-446655440005', 'GT006', 'Exit Gate 1', '550e8400-e29b-41d4-a716-446655440001', true),
('660e8400-e29b-41d4-a716-446655440006', 'GT007', 'Both Gate 1', '550e8400-e29b-41d4-a716-446655440001', true),
-- North Station Gates
('660e8400-e29b-41d4-a716-446655440007', 'GT008', 'North Entry', '550e8400-e29b-41d4-a716-446655440002', true),
('660e8400-e29b-41d4-a716-446655440008', 'GT009', 'North Exit', '550e8400-e29b-41d4-a716-446655440002', true),
('660e8400-e29b-41d4-a716-446655440009', 'GT010', 'Express Gate', '550e8400-e29b-41d4-a716-446655440002', true),
-- South Hub Gates
('660e8400-e29b-41d4-a716-44665544000a', 'GT011', 'South Entry A', '550e8400-e29b-41d4-a716-446655440003', true),
('660e8400-e29b-41d4-a716-44665544000b', 'GT012', 'South Entry B', '550e8400-e29b-41d4-a716-446655440003', true),
('660e8400-e29b-41d4-a716-44665544000c', 'GT013', 'South Exit', '550e8400-e29b-41d4-a716-446655440003', true),
-- East Terminal Gates
('660e8400-e29b-41d4-a716-44665544000d', 'GT014', 'East Main', '550e8400-e29b-41d4-a716-446655440004', true),
('660e8400-e29b-41d4-a716-44665544000e', 'GT015', 'East Secondary', '550e8400-e29b-41d4-a716-446655440004', true),
-- West Station Gates
('660e8400-e29b-41d4-a716-44665544000f', 'GT016', 'West Entry', '550e8400-e29b-41d4-a716-446655440005', true),
('660e8400-e29b-41d4-a716-446655440010', 'GT017', 'West Exit', '550e8400-e29b-41d4-a716-446655440005', true),
-- University Terminal Gates
('660e8400-e29b-41d4-a716-446655440011', 'GT018', 'Student Gate', '550e8400-e29b-41d4-a716-446655440006', true),
('660e8400-e29b-41d4-a716-446655440012', 'GT019', 'Faculty Gate', '550e8400-e29b-41d4-a716-446655440006', true),
('660e8400-e29b-41d4-a716-446655440013', 'GT020', 'Visitor Gate', '550e8400-e29b-41d4-a716-446655440006', true),
-- Shopping Mall Hub Gates
('660e8400-e29b-41d4-a716-446655440014', 'GT021', 'Mall Entry 1', '550e8400-e29b-41d4-a716-446655440007', true),
('660e8400-e29b-41d4-a716-446655440015', 'GT022', 'Mall Entry 2', '550e8400-e29b-41d4-a716-446655440007', true),
('660e8400-e29b-41d4-a716-446655440016', 'GT023', 'Mall Exit', '550e8400-e29b-41d4-a716-446655440007', true),
-- Industrial Zone Gates
('660e8400-e29b-41d4-a716-446655440017', 'GT024', 'Worker Entry', '550e8400-e29b-41d4-a716-446655440008', true),
('660e8400-e29b-41d4-a716-446655440018', 'GT025', 'Cargo Gate', '550e8400-e29b-41d4-a716-446655440008', true),
-- Hospital Station Gates
('660e8400-e29b-41d4-a716-446655440019', 'GT026', 'Emergency Gate', '550e8400-e29b-41d4-a716-446655440009', true),
('660e8400-e29b-41d4-a716-44665544001a', 'GT027', 'General Entry', '550e8400-e29b-41d4-a716-446655440009', true),
('660e8400-e29b-41d4-a716-44665544001b', 'GT028', 'Patient Exit', '550e8400-e29b-41d4-a716-446655440009', true),
-- Sports Complex Gates
('660e8400-e29b-41d4-a716-44665544001c', 'GT029', 'Stadium Entry', '550e8400-e29b-41d4-a716-44665544000a', true),
('660e8400-e29b-41d4-a716-44665544001d', 'GT030', 'VIP Gate', '550e8400-e29b-41d4-a716-44665544000a', true),
-- Beach Terminal Gates
('660e8400-e29b-41d4-a716-44665544001e', 'GT031', 'Beach Access', '550e8400-e29b-41d4-a716-44665544000b', true),
('660e8400-e29b-41d4-a716-44665544001f', 'GT032', 'Tourist Entry', '550e8400-e29b-41d4-a716-44665544000b', true),
-- Mountain Station Gates
('660e8400-e29b-41d4-a716-446655440020', 'GT033', 'Resort Gate', '550e8400-e29b-41d4-a716-44665544000c', true),
('660e8400-e29b-41d4-a716-446655440021', 'GT034', 'Hiking Entry', '550e8400-e29b-41d4-a716-44665544000c', true),
-- Tech Park Hub Gates
('660e8400-e29b-41d4-a716-446655440022', 'GT035', 'Employee Gate', '550e8400-e29b-41d4-a716-44665544000d', true),
('660e8400-e29b-41d4-a716-446655440023', 'GT036', 'Visitor Entry', '550e8400-e29b-41d4-a716-44665544000d', true),
-- Old Town Terminal Gates (inactive terminal)
('660e8400-e29b-41d4-a716-446655440024', 'GT037', 'Heritage Gate', '550e8400-e29b-41d4-a716-44665544000e', false);

INSERT INTO fare_matrix (origin_terminal_id, destination_terminal_id, fare_amount, is_active) VALUES
-- From Central Terminal
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440001', 25.00, true), -- Central to Airport
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440002', 15.00, true), -- Central to North
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440003', 12.00, true), -- Central to South
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440004', 18.00, true), -- Central to East
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440005', 16.00, true), -- Central to West
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440006', 20.00, true), -- Central to University
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440007', 14.00, true), -- Central to Mall
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440008', 22.00, true), -- Central to Industrial
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440009', 17.00, true), -- Central to Hospital
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-44665544000a', 19.00, true), -- Central to Sports
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-44665544000b', 35.00, true), -- Central to Beach
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-44665544000c', 45.00, true), -- Central to Mountain
('550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-44665544000d', 23.00, true), -- Central to Tech Park
-- From Airport Terminal
('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', 25.00, true), -- Airport to Central
('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002', 30.00, true), -- Airport to North
('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440003', 28.00, true), -- Airport to South
('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440006', 32.00, true), -- Airport to University
('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-44665544000b', 40.00, true), -- Airport to Beach
-- From North Station
('550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', 15.00, true), -- North to Central
('550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', 30.00, true), -- North to Airport
('550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440003', 20.00, true), -- North to South
('550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440006', 10.00, true), -- North to University
('550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-44665544000d', 15.00, true), -- North to Tech Park
-- From South Hub
('550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440000', 12.00, true), -- South to Central
('550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440002', 20.00, true), -- South to North
('550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440008', 18.00, true), -- South to Industrial
('550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-44665544000a', 16.00, true), -- South to Sports
-- Additional fare routes for comprehensive testing
('550e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440005', 22.00, true), -- East to West
('550e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440004', 22.00, true), -- West to East
('550e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440007', 8.00, true),  -- University to Mall
('550e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440006', 8.00, true),  -- Mall to University
('550e8400-e29b-41d4-a716-446655440008', '550e8400-e29b-41d4-a716-446655440009', 12.00, true), -- Industrial to Hospital
('550e8400-e29b-41d4-a716-446655440009', '550e8400-e29b-41d4-a716-446655440008', 12.00, true), -- Hospital to Industrial
('550e8400-e29b-41d4-a716-44665544000a', '550e8400-e29b-41d4-a716-44665544000b', 25.00, true), -- Sports to Beach
('550e8400-e29b-41d4-a716-44665544000b', '550e8400-e29b-41d4-a716-44665544000c', 30.00, true), -- Beach to Mountain
('550e8400-e29b-41d4-a716-44665544000c', '550e8400-e29b-41d4-a716-44665544000d', 35.00, true); -- Mountain to Tech Park

INSERT INTO cards (id, card_number, balance, status, issued_date, expiry_date) VALUES
('770e8400-e29b-41d4-a716-446655440000', '1234567890123456', 150.75, 'active', '2023-01-15', '2025-01-15'),
('770e8400-e29b-41d4-a716-446655440001', '2345678901234567', 89.50, 'active', '2023-02-20', '2025-02-20'),
('770e8400-e29b-41d4-a716-446655440002', '3456789012345678', 200.00, 'active', '2023-03-10', '2025-03-10'),
('770e8400-e29b-41d4-a716-446655440003', '4567890123456789', 45.25, 'active', '2023-04-05', '2025-04-05'),
('770e8400-e29b-41d4-a716-446655440004', '5678901234567890', 300.80, 'active', '2023-05-12', '2025-05-12'),
('770e8400-e29b-41d4-a716-446655440005', '6789012345678901', 12.30, 'active', '2023-06-18', '2025-06-18'),
('770e8400-e29b-41d4-a716-446655440006', '7890123456789012', 0.00, 'blocked', '2023-07-22', '2025-07-22'),
('770e8400-e29b-41d4-a716-446655440007', '8901234567890123', 75.60, 'active', '2023-08-30', '2025-08-30'),
('770e8400-e29b-41d4-a716-446655440008', '9012345678901234', 125.40, 'active', '2023-09-14', '2025-09-14'),
('770e8400-e29b-41d4-a716-446655440009', '0123456789012345', 0.00, 'expired', '2022-10-01', '2024-10-01'),
('770e8400-e29b-41d4-a716-44665544000a', '1122334455667788', 95.75, 'active', '2023-11-08', '2025-11-08'),
('770e8400-e29b-41d4-a716-44665544000b', '2233445566778899', 180.20, 'active', '2023-12-03', '2025-12-03'),
('770e8400-e29b-41d4-a716-44665544000c', '3344556677889900', 67.85, 'active', '2024-01-17', '2026-01-17'),
('770e8400-e29b-41d4-a716-44665544000d', '4455667788990011', 0.00, 'blocked', '2024-02-11', '2026-02-11'),
('770e8400-e29b-41d4-a716-44665544000e', '5566778899001122', 220.45, 'active', '2024-03-25', '2026-03-25'),
('770e8400-e29b-41d4-a716-44665544000f', '6677889900112233', 38.90, 'active', '2024-04-09', '2026-04-09'),
('770e8400-e29b-41d4-a716-446655440010', '7788990011223344', 155.30, 'active', '2024-05-14', '2026-05-14'),
('770e8400-e29b-41d4-a716-446655440011', '8899001122334455', 92.15, 'active', '2024-06-07', '2026-06-07'),
('770e8400-e29b-41d4-a716-446655440012', '9900112233445566', 0.00, 'expired', '2022-07-20', '2024-07-20'),
('770e8400-e29b-41d4-a716-446655440013', '0011223344556677', 275.80, 'active', '2024-08-02', '2026-08-02'),
('770e8400-e29b-41d4-a716-446655440014', '1111222233334444', 43.65, 'active', '2024-09-15', '2026-09-15'),
('770e8400-e29b-41d4-a716-446655440015', '2222333344445555', 189.25, 'active', '2024-10-28', '2026-10-28'),
('770e8400-e29b-41d4-a716-446655440016', '3333444455556666', 76.50, 'active', '2024-11-12', '2026-11-12'),
('770e8400-e29b-41d4-a716-446655440017', '4444555566667777', 0.00, 'blocked', '2024-12-05', '2026-12-05'),
('770e8400-e29b-41d4-a716-446655440018', '5555666677778888', 132.90, 'active', '2023-01-30', '2025-01-30'),
('770e8400-e29b-41d4-a716-446655440019', '6666777788889999', 84.40, 'active', '2023-03-18', '2025-03-18'),
('770e8400-e29b-41d4-a716-44665544001a', '7777888899990000', 203.70, 'active', '2023-05-06', '2025-05-06'),
('770e8400-e29b-41d4-a716-44665544001b', '8888999900001111', 51.20, 'active', '2023-07-24', '2025-07-24'),
('770e8400-e29b-41d4-a716-44665544001c', '9999000011112222', 167.35, 'active', '2023-09-11', '2025-09-11'),
('770e8400-e29b-41d4-a716-44665544001d', '0000111122223333', 98.80, 'active', '2023-11-29', '2025-11-29');

INSERT INTO transactions (card_id, gate_id, terminal_id, transaction_type, amount, transaction_time) VALUES
-- Recent transactions for testing
('770e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'tap_in', NULL, '2024-12-29 08:15:00'),
('770e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440001', 'tap_out', 25.00, '2024-12-29 09:30:00'),
('770e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440002', 'tap_in', NULL, '2024-12-29 07:45:00'),
('770e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440014', '550e8400-e29b-41d4-a716-446655440007', 'tap_out', 8.00, '2024-12-29 08:20:00'),
('770e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440006', 'tap_in', NULL, '2024-12-29 09:00:00'),
('770e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-44665544000a', '550e8400-e29b-41d4-a716-446655440003', 'tap_in', NULL, '2024-12-29 10:15:00'),
('770e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440017', '550e8400-e29b-41d4-a716-446655440008', 'tap_out', 18.00, '2024-12-29 11:45:00'),
('770e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440019', '550e8400-e29b-41d4-a716-446655440009', 'tap_in', NULL, '2024-12-29 06:30:00'),
('770e8400-e29b-41d4-a716-446655440005', '660e8400-e29b-41d4-a716-44665544001c', '550e8400-e29b-41d4-a716-44665544000a', 'tap_in', NULL, '2024-12-29 14:00:00'),
('770e8400-e29b-41d4-a716-446655440005', '660e8400-e29b-41d4-a716-44665544001e', '550e8400-e29b-41d4-a716-44665544000b', 'tap_out', 25.00, '2024-12-29 16:30:00'),
-- Historical transactions for data analysis
('770e8400-e29b-41d4-a716-446655440007', '660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', 'tap_in', NULL, '2024-12-28 07:30:00'),
('770e8400-e29b-41d4-a716-446655440007', '660e8400-e29b-41d4-a716-44665544000f', '550e8400-e29b-41d4-a716-446655440005', 'tap_out', 16.00, '2024-12-28 08:45:00'),
('770e8400-e29b-41d4-a716-446655440008', '660e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440001', 'tap_in', NULL, '2024-12-28 12:00:00'),
('770e8400-e29b-41d4-a716-446655440008', '660e8400-e29b-41d4-a716-446655440013', '550e8400-e29b-41d4-a716-446655440006', 'tap_out', 32.00, '2024-12-28 13:15:00'),
('770e8400-e29b-41d4-a716-44665544000a', '660e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-44665544000c', 'tap_in', NULL, '2024-12-28 09:00:00'),
('770e8400-e29b-41d4-a716-44665544000a', '660e8400-e29b-41d4-a716-446655440022', '550e8400-e29b-41d4-a716-44665544000d', 'tap_out', 35.00, '2024-12-28 11:30:00'),
-- Peak hour transactions
('770e8400-e29b-41d4-a716-44665544000b', '660e8400-e29b-41d4-a716-446655440014', '550e8400-e29b-41d4-a716-446655440007', 'tap_in', NULL, '2024-12-27 08:00:00'),
('770e8400-e29b-41d4-a716-44665544000c', '660e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'tap_in', NULL, '2024-12-27 08:05:00'),
('770e8400-e29b-41d4-a716-44665544000d', '660e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440001', 'tap_in', NULL, '2024-12-27 08:10:00'),
('770e8400-e29b-41d4-a716-44665544000e', '660e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440002', 'tap_in', NULL, '2024-12-27 08:15:00'),
('770e8400-e29b-41d4-a716-44665544000f', '660e8400-e29b-41d4-a716-44665544000a', '550e8400-e29b-41d4-a716-446655440003', 'tap_in', NULL, '2024-12-27 08:20:00');

-- Create default admin
INSERT INTO admins (username, password) VALUES
('admin', '$2a$10$i9CyaT6/wYsXdVvvrpjn6Oruj53mpt75.IOt1cvzoGBZk4wbo8lbi');
