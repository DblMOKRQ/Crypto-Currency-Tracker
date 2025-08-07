package models

import "encoding/json"

// TrackedCoin - отслеживаемая криптовалюта
type TrackedCoin struct {
	ID     int64  `json:"id"`
	Symbol string `json:"symbol" validate:"required,alpha"` // Только буквы (BTC, ETH)
}

// CryptoPrice - цена криптовалюты в конкретный момент
type CryptoPrice struct {
	ID        int64   `json:"id"`
	CoinID    int64   `json:"coin_id"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

// AddCoinRequest - запрос на добавление монеты
type AddCoinRequest struct {
	Coin string `json:"coin" validate:"required,alpha"`
}

// GetPriceRequest - запрос на получение цены
type GetPriceRequest struct {
	Coin      string      `json:"coin" validate:"required,alpha"`
	Timestamp json.Number `json:"timestamp" validate:"required"`
}

type GetPriceResponse struct {
	Coin      string  `json:"coin"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

// PriceResponse - ответ с ценой
type PriceResponse struct {
	Coin      string  `json:"coin"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type CoinData struct {
	CurrentPrice float64 `json:"current_price"`
}
