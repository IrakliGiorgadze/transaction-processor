CREATE TABLE user_accounts (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(50) UNIQUE NOT NULL,
    state VARCHAR(10) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    source_type VARCHAR(20) NOT NULL,
    user_account_id INT NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
