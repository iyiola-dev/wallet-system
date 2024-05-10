package transaction

import "time"

type Transaction struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	RequestID string    `json:"request_id" db:"request_id"`
	Amount    int64     `json:"amount" db:"amount"`
	Status    string    `json:"status" db:"status"`
	Type      string    `json:"type" db:"type"`
	Reference string    `json:"reference" db:"reference"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewTransaction creates a new transaction.
func NewTransaction(userID, requestID, reference, transactionType string, amount int64) *Transaction {
	return &Transaction{
		UserID:    userID,
		RequestID: requestID,
		Amount:    amount,
		Status:    "pending",
		Reference: reference,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Type:      transactionType,
	}
}
