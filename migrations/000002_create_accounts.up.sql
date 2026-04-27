CREATE TABLE IF NOT EXISTS accounts (
    account_id      UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_number VARCHAR(50)  NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
