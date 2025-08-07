package scheduler

import (
	"awesomeProject/api"
	"awesomeProject/internal/models"
	"awesomeProject/internal/service"
	"go.uber.org/zap"
	"sync"
	"time"
)

type PricePoller struct {
	coinService  *service.CoinService
	priceUpdates time.Duration
	coinGeckoApi *api.CoinGeckoApi
	log          *zap.Logger
}

func NewPricePoller(coinService *service.CoinService, priceUpdates time.Duration, coinGeckoApi *api.CoinGeckoApi, log *zap.Logger) *PricePoller {
	return &PricePoller{coinService: coinService, priceUpdates: priceUpdates, coinGeckoApi: coinGeckoApi, log: log.Named("PricePoller")}
}

func (p *PricePoller) Start(maxConcurrent int) {
	for {
		coins, err := p.coinService.GetAllCoins()
		if err != nil {
			p.log.Fatal("Error getting all coins", zap.Error(err))
			return
		}
		var wg sync.WaitGroup
		wg.Add(len(coins))

		semaphore := make(chan struct{}, maxConcurrent)
		for i := range coins {
			semaphore <- struct{}{} // Занимаем слот
			go func(coin *models.TrackedCoin) {
				defer wg.Done()
				defer func() { <-semaphore }()
				price, err := p.coinGeckoApi.GetPriceCoin(coin.Symbol)
				if err != nil {
					p.log.Error("Error getting price coin", zap.String("symbol", coin.Symbol), zap.Error(err))
				}
				p.log.Debug("Coin", zap.Int64("ID", coin.ID), zap.String("Symbol", coin.Symbol), zap.Float64("Price", price))
				err = p.coinService.AddNewPrice(&models.CryptoPrice{
					CoinID:    coin.ID,
					Symbol:    coin.Symbol,
					Price:     price,
					Timestamp: time.Now().Unix(),
				})
				if err != nil {
					p.log.Error("Error adding new price", zap.String("symbol", coin.Symbol), zap.Error(err))
				}
			}(coins[i])
			//price, err := p.coinGeckoApi.GetPriceCoin(coin.Symbol)
			//if err != nil {
			//	p.log.Error("Error getting price coin", zap.String("symbol", coin.Symbol), zap.Error(err))
			//}
			//p.log.Debug("Coin", zap.Int64("ID", coin.ID), zap.String("Symbol", coin.Symbol), zap.Float64("Price", price))
			//err = p.coinService.AddNewPrice(&models.CryptoPrice{
			//	CoinID:    coin.ID,
			//	Symbol:    coin.Symbol,
			//	Price:     price,
			//	Timestamp: time.Now().Unix(),
			//})
			//if err != nil {
			//	p.log.Error("Error adding new price", zap.String("symbol", coin.Symbol), zap.Error(err))
			//}

		}

		wg.Wait()
		close(semaphore)

		time.Sleep(p.priceUpdates)
	}
}
