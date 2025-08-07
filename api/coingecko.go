package api

import (
	"awesomeProject/internal/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type CoinGeckoApi struct {
	apiKey     string
	vsCurrency string
	log        *zap.Logger
}

func NewCoinGeckoApi(apiKey string, vsCurrency string, log *zap.Logger) *CoinGeckoApi {
	return &CoinGeckoApi{apiKey: apiKey, vsCurrency: vsCurrency, log: log.Named("CoinGeckoApi")}
}

func (api *CoinGeckoApi) Init() error {
	api.log.Info("Initializing coingecko api")
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/ping?x_cg_demo_api_key=%s", api.apiKey)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode != 200 {
		api.log.Error("Error initializing coingecko api", zap.Int("StatusCode", res.StatusCode))
		return fmt.Errorf("coingecko api init failed with status code %d", res.StatusCode)
	}
	api.log.Info("Successfully initialized coingecko api")
	return nil

}

func (api *CoinGeckoApi) GetPriceCoin(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=%s&symbols=%s&x_cg_demo_api_key=%s", api.vsCurrency, symbol, api.apiKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error getting coins price: status %s", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)

	var coins []models.CoinData
	if err := json.Unmarshal(body, &coins); err != nil {
		return 0, fmt.Errorf("error unmarshalling coins response: %w", err)
	}

	return coins[0].CurrentPrice, nil
}
