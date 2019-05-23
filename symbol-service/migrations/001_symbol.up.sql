CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE symbols(
    id uuid DEFAULT uuid_generate_v4 (),
    name TEXT NOT NULL,
    symbol TEXT NOT NULL,
    market_capitalization NUMERIC,
    sector TEXT,
    industry TEXT,
    exchange TEXT,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX symbols_idx_symbol ON symbols(symbol);
