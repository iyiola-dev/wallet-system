package user

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen --source=repository.go -destination=repository_mock.go -package=user Repository
type Repository interface {
	// GetUserByID  returns the user with the given ID.
	GetUserByID(id string) (*User, error)
}

// service implements the Repository interface.
type service struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

// NewRepository creates a new user repository.
func NewRepository(db *sqlx.DB) Repository {
	return &service{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// GetUserByID returns the user with the given ID.
func (s service) GetUserByID(id string) (*User, error) {
	query, args, err := s.psql.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	if err := s.db.Get(&user, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
