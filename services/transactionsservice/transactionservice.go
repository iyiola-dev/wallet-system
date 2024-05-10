package transactionsservice

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"p-system/repositories/transaction"
	"p-system/services/thirdparty"
	"time"

	"github.com/google/uuid"
)

// TransactionResponse represents the structure of the transaction response
type TransactionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type Request struct {
	Amount float64 `json:"amount" validate:"required"`
	UserID string  `json:"user_id" validate:"required"`
	//type required with one of credit or debit
	Type      string `json:"type" validate:"required,oneof=credit debit"`
	Reference string `json:"reference" validate:"required"`
}

func (s service) HandleTransactionRequest(req Request) (TransactionResponse, error) {

	// Validate if user exists
	user, err := s.userRepo.GetUserByID(req.UserID)
	if err != nil {
		return TransactionResponse{Success: false, Message: "User not found"}, err
	}

	// Validate if user has a wallet
	wallet, err := s.walletRepo.GetWalletByUserID(req.UserID)
	if err != nil {
		return TransactionResponse{Success: false, Message: "Wallet not found"}, err
	}

	// Convert balance to float64 by dividing by 100
	balance := float64(wallet.Balance) / 100

	// Convert amount to int64 by multiplying by 100
	amount := req.Amount * 100

	// If type is debit, check if user has enough balance
	if req.Type == "debit" {
		if balance < req.Amount {
			return TransactionResponse{Success: false, Message: "Insufficient balance"}, nil
		}
	}

	requestID := uuid.NewString()

	// Create transaction
	transaction := transaction.NewTransaction(req.UserID, requestID, req.Reference, req.Type, int64(amount))

	// Create transaction
	transaction, err = s.transactionRepo.Create(transaction)
	if err != nil {
		log.Println("error", err)
		return TransactionResponse{Success: false, Message: "Failed to create transaction"}, err
	}

	// Send request to third party to make payment with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Make payment
	_, err = s.thirdPartyService.MakePayment(thirdparty.Transaction{
		AccountID: user.ID,
		Reference: req.Reference,
		Amount:    req.Amount,
	}, ctx)

	if err != nil {
		log.Println("error", err)
		_, newErr := s.transactionRepo.UpdateTransactionToFailed(transaction.ID)
		if newErr != nil {
			log.Println("error", err)
			return TransactionResponse{Success: false, Message: "Failed to update transaction"}, newErr
		}

		return TransactionResponse{Success: false, Message: "Failed to make payment"}, err
	}

	// Update wallet
	if req.Type == "debit" {
		// Call debit wallet
		wallet, err = s.walletRepo.DebitWallet(wallet, *transaction, int64(amount))

		if err != nil {
			log.Println("error", err)
			return TransactionResponse{Success: false, Message: "Failed to debit wallet"}, err
		}

	} else if req.Type == "credit" {
		// Call credit wallet
		wallet, err = s.walletRepo.CreditWallet(wallet, *transaction, int64(amount))

		if err != nil {
			log.Println("error", err)
			return TransactionResponse{Success: false, Message: "Failed to credit wallet"}, err
		}

	}

	return TransactionResponse{Success: true, Message: "Transaction successful"}, nil
}

func (s service) HandleTransaction(w http.ResponseWriter, r *http.Request) {

	var req Request

	//decode request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return

	}

	//call service method
	resp, err := s.HandleTransactionRequest(req)

	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, resp)
		return
	}

	sendJSONResponse(w, http.StatusOK, resp)

}

// sendJSONResponse sends a JSON response with the specified status code and data
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
