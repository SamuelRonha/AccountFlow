
CREATE TABLE IF NOT EXISTS transactions (
    transaction_id    UUID          PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id        UUID          NOT NULL REFERENCES accounts(account_id) ON DELETE CASCADE,
    operation_type_id INT           NOT NULL REFERENCES operation_types(operation_type_id),
    amount            NUMERIC(15,2) NOT NULL,
    event_date        TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
