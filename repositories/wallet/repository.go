package wallet

import (
	"fmt"
	"log"
	"p-system/repositories/transaction"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen --source=repository.go -destination=respository_mock.go -package=wallet Repository
type Repository interface {
	// Create creates a new wallet.
	Create(*Wallet) (*Wallet, error)
	// GetWalletByUserID returns the wallet with the given user id.
	GetWalletByUserID(string) (*Wallet, error)
	// CreditWallet updates the balance of a wallet.
	CreditWallet(*Wallet, transaction.Transaction, int64) (*Wallet, error)
	// DebitWallet updates the balance of a wallet.
	DebitWallet(*Wallet, transaction.Transaction, int64) (*Wallet, error)
}

type service struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

func NewRepository(db *sqlx.DB) Repository {
	return &service{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create creates a new wallet.
func (s service) Create(wallet *Wallet) (*Wallet, error) {
	query, args, err := s.psql.Insert("wallets").
		Columns("user_id", "balance", "created_at", "updated_at").
		Values(wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	var w Wallet
	if err := s.db.Get(&w, query, args...); err != nil {
		return nil, err
	}

	return &w, nil
}

// GetWalletByUserID returns the wallet with the given user id.
func (s service) GetWalletByUserID(userID string) (*Wallet, error) {
	query, args, err := s.psql.Select("*").
		From("wallets").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	log.Println("query", Log(query, args...))

	var w Wallet
	if err := s.db.Get(&w, query, args...); err != nil {
		log.Println("error", err)
		return nil, err
	}

	return &w, nil
}

// CreditWallet updates the balance of a wallet.
func (s service) CreditWallet(wallet *Wallet, transaction transaction.Transaction, amount int64) (*Wallet, error) {
	//use transaction to ensure atomicity
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}

	//update transaction status to completed
	_, err = tx.Exec("UPDATE transactions SET status = $1,updated_at = $2 WHERE id = $3", "completed", time.Now(), transaction.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//get the transaction id
	wallet.TransactionID = &transaction.ID

	//lock the wallet row to prevent concurrent updates
	_, err = tx.Exec("SELECT * FROM wallets WHERE id = $1 FOR UPDATE", wallet.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//update the wallet and transaction id with the values passed

	_, err = tx.Exec("UPDATE wallets SET balance = balance + $1, updated_at = $2, transaction_id = $3 WHERE id = $4", amount, time.Now(), wallet.TransactionID, wallet.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	//fetch the updated wallet
	query, args, err := s.psql.Select("*").
		From("wallets").
		Where(sq.Eq{"id": wallet.ID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var w Wallet

	if err := s.db.Get(&w, query, args...); err != nil {
		return nil, err
	}

	return &w, nil

}

// DebitWallet updates the balance of a wallet.
func (s service) DebitWallet(wallet *Wallet, transaction transaction.Transaction, amount int64) (*Wallet, error) {
	//use transaction to ensure atomicity
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}

	//update transaction status to completed
	_, err = tx.Exec("UPDATE transactions SET status = $1,updated_at = $2 WHERE id = $3", "completed", time.Now(), transaction.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	wallet.TransactionID = &transaction.ID

	//lock the wallet row to prevent concurrent updates
	_, err = tx.Exec("SELECT * FROM wallets WHERE id = $1 FOR UPDATE", wallet.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//update the wallet and transaction id with the values passed

	_, err = tx.Exec("UPDATE wallets SET balance = balance - $1, updated_at = $2, transaction_id = $3 WHERE id = $4", amount, time.Now(), wallet.TransactionID, wallet.ID)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	//fetch the updated wallet
	query, args, err := s.psql.Select("*").
		From("wallets").
		Where(sq.Eq{"id": wallet.ID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var w Wallet

	if err := s.db.Get(&w, query, args...); err != nil {
		return nil, err
	}

	return &w, nil
}

// Log provides a pretty print version of the query and parameters.
func Log(query string, args ...interface{}) string {
	for i, arg := range args {
		n := fmt.Sprintf("$%d", i+1)

		var a string
		switch v := arg.(type) {
		case string:
			a = fmt.Sprintf("%q", v)
		case []byte:
			a = string(v)
		case []string:
			a = strings.Join(v, ",")
		default:
			a = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, n, a, 1)
	}

	return query
}
