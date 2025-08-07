-- +goose Up
-- Создаем таблицу отслеживаемых монет
CREATE TABLE tracked_coins (
                               id SERIAL PRIMARY KEY,
                               symbol VARCHAR(10) UNIQUE NOT NULL,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу цен
CREATE TABLE coin_prices (
                             id SERIAL PRIMARY KEY,
                             coin_id INTEGER NOT NULL REFERENCES tracked_coins(id) ON DELETE CASCADE,
                             price DECIMAL(20,8) NOT NULL,
                             timestamp BIGINT NOT NULL,
                             UNIQUE(coin_id, timestamp)
);

-- Индексы для ускорения запросов
CREATE INDEX idx_coin_prices_coin_id ON coin_prices(coin_id);
CREATE INDEX idx_coin_prices_timestamp ON coin_prices(timestamp);

-- +goose Down
DROP TABLE IF EXISTS coin_prices;
DROP TABLE IF EXISTS tracked_coins;