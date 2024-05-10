package thirdparty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)


type service struct {

}

//go:generate mockgen --source=thirdparty.go -destination=thirdparty_mock.go -package=thirdparty Service
type Service interface {
	GetTransaction(reference, accountID string) (*Transaction, error)
	MakePayment(req Transaction, ctx context.Context) (*Transaction, error)

}

func NewService() Service {
	return &service{}
}

type Transaction struct {
	AccountID string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}

func(s *service) GetTransaction(reference, accountID string) (*Transaction, error) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulating response for testing
		mockTransaction := Transaction{
			AccountID: accountID,
			Reference: reference, // Use the provided reference
			Amount:    100.0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockTransaction)
	}))
	defer server.Close()

	// Send a GET request to the mock server
	resp, err := http.Get(server.URL + "/third-party/payments/" + reference)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response body into a Transaction object
	var transaction Transaction
	err = json.NewDecoder(resp.Body).Decode(&transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// MakePayment sends a payment request to a third-party service.
func(s *service) MakePayment(req Transaction, ctx context.Context) (*Transaction, error) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulating response for testing
		mockTransaction := Transaction{
			AccountID: req.AccountID,
			Reference: req.Reference,
			Amount:    req.Amount,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockTransaction)
	}))
	defer server.Close()

	// Create a new request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Send a POST request to the mock server
	resp, err := http.Post(server.URL+"/third-party/payments", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response body into a Transaction object
	var transaction Transaction
	err = json.NewDecoder(resp.Body).Decode(&transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
