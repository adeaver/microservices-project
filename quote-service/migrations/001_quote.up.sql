CREATE TABLE quotes(
    symbol TEXT NOT NULL,
    time TIMESTAMP NOT NULL,
    open_price_cents NUMERIC,
    high_price_cents NUMERIC,
    low_price_cents NUMERIC,
    close_price_cents NUMERIC,
    volume_shares NUMERIC
);

CREATE UNIQUE INDEX quotes_symbol_time_idx ON quotes (symbol, time);
