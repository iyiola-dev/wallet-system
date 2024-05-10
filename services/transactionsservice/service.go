package transactionsservice

import (
	"net/http"
	"p-system/repositories/transaction"
	"p-system/repositories/user"
	"p-system/repositories/wallet"
	"p-system/services/thirdparty"
)

type service struct {
	userRepo          user.Repository
	transactionRepo   transaction.Repository
	walletRepo        wallet.Repository
	thirdPartyService thirdparty.Service
}

type Service interface {
	HandleTransaction(w http.ResponseWriter, r *http.Request)
}

func NewService(userRepo user.Repository, transactionRepo transaction.Repository, walletRepo wallet.Repository, thirdpartyService thirdparty.Service) Service {
	return &service{
		userRepo:          userRepo,
		transactionRepo:   transactionRepo,
		walletRepo:        walletRepo,
		thirdPartyService: thirdpartyService,
	}
}
