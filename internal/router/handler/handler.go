package handler

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/service"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	coinService *service.CoinService
}

func NewHandler(coinService *service.CoinService) *Handler {
	return &Handler{coinService: coinService}
}

func (h *Handler) AddCoin(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value("logger").(*zap.Logger)
	if r.Method != http.MethodPost {
		log.Warn("Invalid request method", zap.String("path", r.URL.Path), zap.String("method", r.Method))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("Handling add coin")
	// Извлечение данных из запроса
	var addReq models.AddCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&addReq); err != nil {
		log.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if addReq.Coin == "" {
		log.Warn("Coin is required")
		http.Error(w, "Coin is required", http.StatusBadRequest)
		return
	}

	if err := h.coinService.AddCoin(&addReq); err != nil {
		log.Warn("Failed to add coin", zap.Error(err))
		http.Error(w, "Failed to add coin", http.StatusInternalServerError)
		return
	}

	log.Info("Added coin", zap.String("coin", addReq.Coin))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetCoin(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value("logger").(*zap.Logger)
	if r.Method != http.MethodPost {
		log.Warn("Invalid request method", zap.String("path", r.URL.Path), zap.String("method", r.Method))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("Handling get coin", zap.String("path", r.URL.Path), zap.String("method", r.Method))

	// Извлечение данных из запроса
	var getReq models.GetPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&getReq); err != nil {
		log.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if getReq.Coin == "" {
		log.Warn("Coin is required")
		http.Error(w, "Coin is required", http.StatusBadRequest)
		return
	}
	checkTimestamp, err := getReq.Timestamp.Int64()
	if err != nil {
		log.Warn("Invalid timestamp", zap.Error(err))
		http.Error(w, "Invalid timestamp", http.StatusBadRequest)
		return
	}
	if checkTimestamp <= 0 {
		log.Warn("Timestamp must be positive")
		http.Error(w, "Timestamp must be positive", http.StatusBadRequest)
		return
	}

	price, err := h.coinService.GetPrice(&getReq)
	if err != nil {
		log.Warn("Failed to get coin", zap.Error(err))
		http.Error(w, "Failed to get coin", http.StatusInternalServerError)
		return
	}
	log.Info("Get coin", zap.String("coin", getReq.Coin), zap.Float64("price", price.Price))

	response := models.GetPriceResponse{
		Coin:      price.Symbol,
		Price:     price.Price,
		Timestamp: price.Timestamp,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Warn("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteCoin(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value("logger").(*zap.Logger)
	if r.Method != http.MethodPost {
		log.Warn("Invalid request method", zap.String("path", r.URL.Path))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	log.Info("Handling delete coin", zap.String("path", r.URL.Path))

	// Извлечение данных из запроса
	var deleteReq models.AddCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&deleteReq); err != nil {
		log.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if deleteReq.Coin == "" {
		log.Warn("Coin is required")
		http.Error(w, "Coin is required", http.StatusBadRequest)
		return
	}

	if err := h.coinService.RemoveCoin(&deleteReq); err != nil {
		log.Warn("Failed to remove coin", zap.Error(err))
		http.Error(w, "Failed to remove coin", http.StatusInternalServerError)
		return
	}
	log.Info("Removed coin", zap.String("coin", deleteReq.Coin))
	w.WriteHeader(http.StatusOK)

}
