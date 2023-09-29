CREATE TABLE accounts (
    account_id     UUID         NOT NULL,
    "name"         VARCHAR(255) NOT NULL,
    account_number VARCHAR(255) NOT NULL,
    customer_id    UUID         NOT NULL,
    PRIMARY KEY (account_id)
);

CREATE TABLE account_logs (
    account_log_id UUID      NOT NULL,
    account        JSONB     NOT NULL,
    "time"         TIMESTAMP NOT NULL,
    account_id     UUID      GENERATED ALWAYS AS ((account->>'account_id')::UUID) STORED,
    customer_id    UUID      GENERATED ALWAYS AS ((account->>'customer_id')::UUID) STORED,
    PRIMARY KEY (account_log_id)
);

CREATE TABLE transaction_logs (
    transaction_log_id UUID      NOT NULL,
    "transaction"      JSONB     NOT NULL,
    "time"             TIMESTAMP NOT NULL,
    account_id         UUID      GENERATED ALWAYS AS (("transaction"->>'account_id')::UUID) STORED,
    PRIMARY KEY (transaction_log_id)
);

CREATE TABLE accounts_balance (
    account_id UUID          NOT NULL,
    balance    DECIMAL(10,2) NOT NULL,
    PRIMARY KEY (account_id)
);
