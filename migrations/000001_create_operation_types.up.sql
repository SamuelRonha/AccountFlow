CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Operation types are seeded once and never change at runtime.
-- 1=Normal Purchase, 2=Purchase with Installments, 3=Withdrawal, 4=Credit Voucher
CREATE TABLE IF NOT EXISTS operation_types (
    operation_type_id SERIAL      PRIMARY KEY,
    description       VARCHAR(100) NOT NULL
);

INSERT INTO operation_types (operation_type_id, description) VALUES
    (1, 'Normal Purchase'),
    (2, 'Purchase with Installments'),
    (3, 'Withdrawal'),
    (4, 'Credit Voucher')
ON CONFLICT DO NOTHING;
