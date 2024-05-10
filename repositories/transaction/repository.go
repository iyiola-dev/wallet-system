package transaction

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)


//go:generate mockgen --source=repository.go -destination=respository_mock.go -package=transaction Repository
type Repository interface {
	// Create creates a new transaction.
	Create(*Transaction) (*Transaction, error)
	// GetTransactionByReference returns the transaction with the given reference.
	GetTransactionByReference(string) (*Transaction, error)

	//UpdateTransactionToFailed updates a transaction
	UpdateTransactionToFailed(id string) (*Transaction, error)
}

// service implements the Repository interface.
type service struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}



// NewRepository creates a new transaction repository.
func NewRepository(db *sqlx.DB) Repository {
	return &service{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create creates a new transaction.
func (s service) Create(transaction *Transaction) (*Transaction, error) {
	query, args, err := s.psql.Insert("transactions").
		Columns("user_id", "request_id", "type", "amount", "status", "reference", "created_at", "updated_at").
		Values(transaction.UserID, transaction.RequestID, transaction.Type, transaction.Amount, transaction.Status, transaction.Reference, transaction.CreatedAt, transaction.UpdatedAt).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	var t Transaction
	if err := s.db.Get(&t, query, args...); err != nil {
		//check if error is due to constraint violation
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return nil, errors.New("transaction already exists")
			}
		}

		return nil, err
	}

	return &t, nil
}

// GetTransactionByReference returns the transaction with the given reference.
func (s service) GetTransactionByReference(reference string) (*Transaction, error) {
	query, args, err := s.psql.Select("*").
		From("transactions").
		Where(sq.Eq{"reference": reference}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var t Transaction
	if err := s.db.Get(&t, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return &t, nil
}

func (s service) UpdateTransactionToFailed(id string) (*Transaction, error) {
	query, args, err := s.psql.Update("transactions").
		Set("status", "failed").
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	var t Transaction
	if err := s.db.Get(&t, query, args...); err != nil {
		return nil, err
	}

	return &t, nil
}
