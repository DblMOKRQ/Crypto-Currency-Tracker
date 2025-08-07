package repository

import (
	"awesomeProject/internal/models"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type Repository struct {
	db  *sql.DB
	log *zap.Logger
}

func (s *Storage) NewRepository() *Repository {
	return &Repository{db: s.db, log: s.log.Named("Repository")}
}
func (r *Repository) AddCoin(coin *models.TrackedCoin) error {
	stmt, err := r.db.Prepare(`
        INSERT INTO tracked_coins (symbol) 
        VALUES ($1) 
        ON CONFLICT (symbol) DO NOTHING`)
	if err != nil {
		r.log.Error("Failed to insert coin", zap.Error(err))
		return fmt.Errorf("failed to insert coin: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(coin.Symbol)
	return err
}
func (r *Repository) RemoveCoin(coin *models.TrackedCoin) error {
	stmt, err := r.db.Prepare(`
		DELETE FROM tracked_coins 
        WHERE symbol = $1`)
	if err != nil {
		r.log.Error("Failed to remove coin", zap.Error(err))
		return fmt.Errorf("failed to remove coin: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(coin.Symbol)
	return err
}
func (r *Repository) GetPrice(coin *models.GetPriceRequest) (*models.CryptoPrice, error) {
	var price models.CryptoPrice
	query := `
        SELECT 
            cp.id,
            cp.coin_id,
            tc.symbol,
            cp.price,
            cp.timestamp
        FROM coin_prices cp
        JOIN tracked_coins tc ON tc.id = cp.coin_id
        WHERE tc.symbol = $1
        ORDER BY ABS(cp.timestamp - $2)
        LIMIT 1`

	err := r.db.QueryRow(query, coin.Coin, coin.Timestamp).Scan(
		&price.ID,
		&price.CoinID,
		&price.Symbol,
		&price.Price,
		&price.Timestamp,
	)

	if err != nil {
		r.log.Error("Failed to get price", zap.Error(err), zap.String("coin", coin.Coin))
		return nil, fmt.Errorf("failed to get price: %w", err)
	}
	return &price, nil
}
func (r *Repository) GetAllCoins() ([]*models.TrackedCoin, error) {
	query := `SELECT id, symbol FROM tracked_coins ORDER BY symbol`
	r.log.Debug("Getting all coins")
	rows, err := r.db.Query(query)
	if err != nil {
		r.log.Error("Failed to get all coins", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch tracked coins: %w", err)
	}
	defer rows.Close()

	var coins []*models.TrackedCoin
	for rows.Next() {
		var coin models.TrackedCoin
		if err := rows.Scan(
			&coin.ID,
			&coin.Symbol,
		); err != nil {
			r.log.Error("Failed to scan coin", zap.Error(err))
			return nil, fmt.Errorf("failed to scan coin: %w", err)
		}
		coins = append(coins, &coin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	r.log.Debug("Successfully fetched all coins", zap.Int("count", len(coins)))

	return coins, nil
}

func (r *Repository) AddNewPrice(coin *models.CryptoPrice) error {
	r.log.Debug("Adding new price", zap.String("symbol", coin.Symbol))

	// Используем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Safe to call if tx is already committed

	// Проверяем существование монеты
	var exists bool
	err = tx.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM tracked_coins WHERE id = $1)",
		coin.CoinID,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check coin existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("coin with id %d not found", coin.CoinID)
	}

	// Вставляем цену
	_, err = tx.Exec(`
        INSERT INTO coin_prices (coin_id, price, timestamp)
        VALUES ($1, $2, $3)
        ON CONFLICT (coin_id, timestamp) DO UPDATE
        SET price = EXCLUDED.price`,
		coin.CoinID,
		coin.Price,
		coin.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to insert price: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.log.Debug("Successfully added new price",
		zap.String("symbol", coin.Symbol),
		zap.Float64("price", coin.Price))
	return nil
}
