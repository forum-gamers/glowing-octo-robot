CREATE TABLE Transaction (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    userId UUID NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    type VARCHAR(100) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'failed','cancel','refund','settlement','deny','expire')),
    transactionDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT NULL,
    detail TEXT NULL,
    discount DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
