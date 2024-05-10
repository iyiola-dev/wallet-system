package wallet

import "time"

type Wallet struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	Balance       int64     `json:"balance" db:"balance"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	TransactionID *string   `json:"transaction_id" db:"transaction_id"`
}

// NewWallet creates a new wallet.
func NewWallet(userID string, balance int64) *Wallet {
	return &Wallet{
		UserID:    userID,
		Balance:   balance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
