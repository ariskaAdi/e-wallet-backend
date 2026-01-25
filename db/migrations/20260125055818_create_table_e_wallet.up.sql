CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    public_id VARCHAR(100) NOT NULL,
    otp VARCHAR(4) NOT NULL,
    password VARCHAR(100) NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE wallet (
    id SERIAL PRIMARY KEY,
    wallet_public_id VARCHAR(100) NOT NULL,
    user_public_id VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    balance BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(100) NOT NULL,
    wallet_public_id VARCHAR(100) NOT NULL,
    sof_number VARCHAR(100) NOT NULL,
    dof_number VARCHAR(100) NOT NULL,
    type VARCHAR(100) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(100) NOT NULL,
    reference VARCHAR(100) NOT NULL,
    description VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE transfer_inquiry (
    inquiry_key VARCHAR(100) NOT NULL PRIMARY KEY,
    sof_number VARCHAR(100) NOT NULL,
    dof_number VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expired_at TIMESTAMP NOT NULL
);

CREATE TABLE topup (
    id SERIAL PRIMARY KEY,
    topup_id VARCHAR(100) NOT NULL,
    user_public_id VARCHAR(100) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(100) NOT NULL,
    snap_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

