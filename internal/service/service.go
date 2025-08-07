package service

import (
	"awesomeProject/internal/models"
	"errors"
	"regexp"
	"strings"
)

type repository interface {
	AddCoin(coin *models.TrackedCoin) error
	RemoveCoin(coin *models.TrackedCoin) error
	GetPrice(coin *models.GetPriceRequest) (*models.CryptoPrice, error)
	GetAllCoins() ([]*models.TrackedCoin, error)
	AddNewPrice(coin *models.CryptoPrice) error
}

type CoinService struct {
	repo repository
}

func NewCoinService(repo repository) *CoinService {
	return &CoinService{repo: repo}
}

func (c *CoinService) AddCoin(req *models.AddCoinRequest) error {
	if !validateSymbol(req.Coin) {
		return errors.New("invalid coin")
	}
	coin := models.TrackedCoin{
		Symbol: strings.ToUpper(req.Coin),
	}
	return c.repo.AddCoin(&coin)
}
func (c *CoinService) RemoveCoin(req *models.AddCoinRequest) error {
	if !validateSymbol(req.Coin) {
		return errors.New("invalid coin")
	}
	coin := models.TrackedCoin{
		Symbol: strings.ToUpper(req.Coin),
	}
	return c.repo.RemoveCoin(&coin)
}
func (c *CoinService) GetPrice(coin *models.GetPriceRequest) (*models.CryptoPrice, error) {
	if !validateSymbol(coin.Coin) {
		return nil, errors.New("invalid coin")
	}
	coin.Coin = strings.ToUpper(coin.Coin)
	return c.repo.GetPrice(coin)
}
func (c *CoinService) GetAllCoins() ([]*models.TrackedCoin, error) {
	return c.repo.GetAllCoins()
}
func (c *CoinService) AddNewPrice(coin *models.CryptoPrice) error {
	if !validateSymbol(coin.Symbol) {
		return errors.New("invalid symbol")
	}
	if coin.Price <= 0 {
		return errors.New("invalid price")
	}
	if coin.Timestamp == 0 {
		return errors.New("invalid timestamp")
	}
	return c.repo.AddNewPrice(coin)
}

func validateSymbol(symbol string) bool {
	matched, _ := regexp.MatchString(`^[A-Za-z]{1,10}$`, symbol)
	return matched
}
