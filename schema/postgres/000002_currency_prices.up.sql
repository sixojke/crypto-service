CREATE TABLE currency_prices (
    symbol VARCHAR(16) NOT NULL,
    price DECIMAL NOT NULL,
    timestamp TIMESTAMP,
    PRIMARY KEY (symbol, timestamp)
);

CREATE INDEX idx_symbol_timestamp ON currency_prices (symbol, timestamp DESC);