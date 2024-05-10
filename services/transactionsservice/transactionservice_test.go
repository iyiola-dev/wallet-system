package transactionsservice

import (
	"p-system/repositories/transaction"
	"p-system/repositories/user"
	"p-system/repositories/wallet"
	"p-system/services/thirdparty"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDebitHandleTransactionRequest(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories and third party
	mockUserRepo := user.NewMockRepository(ctrl)
	mockWalletRepo := wallet.NewMockRepository(ctrl)
	mockTransactionRepo := transaction.NewMockRepository(ctrl)
	mockThirdParty := thirdparty.NewMockService(ctrl)

	// Create a sample request
	req := Request{
		Amount: 100.0,
		UserID: "user123",
		Type:   "debit",
	}

	// Create a sample user
	mockUser := user.User{
		ID: "user123",
	}

	// Create a sample wallet
	mockWallet := wallet.Wallet{
		UserID:  "user123",
		Balance: 20000, // $200.00 in cents
	}

	// Create a sample transaction
	mockTransaction := transaction.Transaction{
		UserID:    "user123",
		RequestID: uuid.NewString(),
		Reference: "ref123",
		Type:      "debit",
		Amount:    10000, // $100.00 in cents
	}

	mockthirdPartyTransaction := thirdparty.Transaction{
		AccountID: "user123",
		Reference: "ref123",
		Amount:    100.0,
	}

	// Set up expectations
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(&mockUser, nil)
	mockWalletRepo.EXPECT().GetWalletByUserID(req.UserID).Return(&mockWallet, nil)
	mockTransactionRepo.EXPECT().Create(gomock.Any()).Return(&mockTransaction, nil)
	mockWalletRepo.EXPECT().DebitWallet(&mockWallet, gomock.Any(), mockTransaction.Amount).Return(&mockWallet, nil)
	mockThirdParty.EXPECT().MakePayment(gomock.Any(), gomock.Any()).Return(&mockthirdPartyTransaction, nil)

	// Create the service with mocked dependencies
	svc := service{
		userRepo:          mockUserRepo,
		walletRepo:        mockWalletRepo,
		transactionRepo:   mockTransactionRepo,
		thirdPartyService: mockThirdParty,
	}

	// Call the method
	resp, err := svc.HandleTransactionRequest(req)

	// Check the result
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Transaction successful", resp.Message)
}

func TestCreditHandleTransactionRequest(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories and third party
	mockUserRepo := user.NewMockRepository(ctrl)
	mockWalletRepo := wallet.NewMockRepository(ctrl)
	mockTransactionRepo := transaction.NewMockRepository(ctrl)
	mockThirdParty := thirdparty.NewMockService(ctrl)

	// Create a sample request
	req := Request{
		Amount: 100.0,
		UserID: "user123",
		Type:   "credit",
	}

	// Create a sample user
	mockUser := user.User{
		ID: "user123",
	}

	// Create a sample wallet
	mockWallet := wallet.Wallet{
		UserID:  "user123",
		Balance: 20000, // $200.00 in cents
	}

	// Create a sample transaction
	mockTransaction := transaction.Transaction{
		UserID:    "user123",
		RequestID: uuid.NewString(),
		Reference: "ref123",
		Type:      "credit",
		Amount:    10000, // $100.00 in cents
	}

	mockthirdPartyTransaction := thirdparty.Transaction{
		AccountID: "user123",
		Reference: "ref123",
		Amount:    100.0,
	}

	// Set up expectations
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(&mockUser, nil)
	mockWalletRepo.EXPECT().GetWalletByUserID(req.UserID).Return(&mockWallet, nil)
	mockTransactionRepo.EXPECT().Create(gomock.Any()).Return(&mockTransaction, nil)
	mockWalletRepo.EXPECT().CreditWallet(&mockWallet, gomock.Any(), mockTransaction.Amount).Return(&mockWallet, nil)
	mockThirdParty.EXPECT().MakePayment(gomock.Any(), gomock.Any()).Return(&mockthirdPartyTransaction, nil)

	// Create the service with mocked dependencies
	svc := service{
		userRepo:          mockUserRepo,
		walletRepo:        mockWalletRepo,
		transactionRepo:   mockTransactionRepo,
		thirdPartyService: mockThirdParty,
	}

	// Call the method
	resp, err := svc.HandleTransactionRequest(req)

	// Check the result
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Transaction successful", resp.Message)
}

func TestHandleTransactionRequest_UserNotFound(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories and third party
	mockUserRepo := user.NewMockRepository(ctrl)

	// Create a sample request
	req := Request{
		Amount: 100.0,
		UserID: "user123",
		Type:   "debit",
	}

	// Set up expectations
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(nil, assert.AnError)

	// Create the service with mocked dependencies
	svc := service{
		userRepo: mockUserRepo,
	}

	// Call the method
	resp, err := svc.HandleTransactionRequest(req)

	// Check the result
	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "User not found", resp.Message)
}

func TestHandleTransactionRequest_WalletNotFound(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories and third party
	mockUserRepo := user.NewMockRepository(ctrl)
	mockWalletRepo := wallet.NewMockRepository(ctrl)

	// Create a sample request
	req := Request{
		Amount: 100.0,
		UserID: "user123",
		Type:   "debit",
	}

	// Create a sample user
	mockUser := user.User{
		ID: "user123",
	}

	// Set up expectations
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(&mockUser, nil)
	mockWalletRepo.EXPECT().GetWalletByUserID(req.UserID).Return(nil, assert.AnError)

	// Create the service with mocked dependencies
	svc := service{
		userRepo:   mockUserRepo,
		walletRepo: mockWalletRepo,
	}

	// Call the method
	resp, err := svc.HandleTransactionRequest(req)

	// Check the result
	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Wallet not found", resp.Message)
}

func TestHandleTransactionRequest_InsufficientBalance(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories and third party
	mockUserRepo := user.NewMockRepository(ctrl)
	mockWalletRepo := wallet.NewMockRepository(ctrl)

	// Create a sample request
	req := Request{
		Amount: 100.0,
		UserID: "user123",
		Type:   "debit",
	}

	// Create a sample user
	mockUser := user.User{
		ID: "user123",
	}

	// Create a sample wallet
	mockWallet := wallet.Wallet{
		UserID:  "user123",
		Balance: 5000, // $50.00 in cents
	}

	// Set up expectations
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(&mockUser, nil)
	mockWalletRepo.EXPECT().GetWalletByUserID(req.UserID).Return(&mockWallet, nil)

	// Create the service with mocked dependencies
	svc := service{
		userRepo:   mockUserRepo,
		walletRepo: mockWalletRepo,
	}

	// Call the method
	resp, err := svc.HandleTransactionRequest(req)

	// Check the result
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Insufficient balance", resp.Message)
}

func TestHandleTransactionRequest_FailedToCreateTransaction(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories and third party
	mockUserRepo := user.NewMockRepository(ctrl)
	mockWalletRepo := wallet.NewMockRepository(ctrl)
	mockTransactionRepo := transaction.NewMockRepository(ctrl)

	// Create a sample request
	req := Request{
		Amount: 100.0,
		UserID: "user123",
		Type:   "debit",
	}

	// Create a sample user
	mockUser := user.User{
		ID: "user123",
	}

	// Create a sample wallet
	mockWallet := wallet.Wallet{
		UserID:  "user123",
		Balance: 20000, // $200.00 in cents
	}

	// Set up expectations
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(&mockUser, nil)
	mockWalletRepo.EXPECT().GetWalletByUserID(req.UserID).Return(&mockWallet, nil)
	mockTransactionRepo.EXPECT().Create(gomock.Any()).Return(nil, assert.AnError)

	// Create the service with mocked dependencies
	svc := service{
		userRepo:        mockUserRepo,
		walletRepo:      mockWalletRepo,
		transactionRepo: mockTransactionRepo,
	}

	// Call the method
	resp, err := svc.HandleTransactionRequest(req)

	// Check the result
	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Failed to create transaction", resp.Message)
}

func TestHandleTransactionRequest_FailedToDebitWallet(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories and third party
	mockUserRepo := user.NewMockRepository(ctrl)
	mockWalletRepo := wallet.NewMockRepository(ctrl)
	mockTransactionRepo := transaction.NewMockRepository(ctrl)
	mockThirdPartyRepo := thirdparty.NewMockService(ctrl)

	// Create a sample request
	req := Request{
		Amount: 100.0,
		UserID: "user123",
		Type:   "debit",
	}

	// Create a sample user
	mockUser := user.User{
		ID: "user123",
	}

	// Create a sample wallet
	mockWallet := wallet.Wallet{
		UserID:  "user123",
		Balance: 20000, // $200.00 in cents
	}

	// Create a sample transaction
	mockTransaction := transaction.Transaction{
		UserID:    "user123",
		RequestID: uuid.NewString(),
		Reference: "ref123",
		Type:      "debit",
		Amount:    10000, // $100.00 in cents
	}

	// Set up expectations
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(&mockUser, nil)
	mockWalletRepo.EXPECT().GetWalletByUserID(req.UserID).Return(&mockWallet, nil)
	mockTransactionRepo.EXPECT().Create(gomock.Any()).Return(&mockTransaction, nil)
	mockThirdPartyRepo.EXPECT().MakePayment(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockWalletRepo.EXPECT().DebitWallet(&mockWallet, gomock.Any(), mockTransaction.Amount).Return(nil, assert.AnError)

	// Create the service with mocked dependencies
	svc := service{
		userRepo:          mockUserRepo,
		walletRepo:        mockWalletRepo,
		transactionRepo:   mockTransactionRepo,
		thirdPartyService: mockThirdPartyRepo,
	}

	// Call the method
	resp, err := svc.HandleTransactionRequest(req)

	// Check the result
	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Failed to debit wallet", resp.Message)
}

func TestHandleTransactionRequest_FailedToMakePayment(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repositories and third party
	mockUserRepo := user.NewMockRepository(ctrl)
	mockWalletRepo := wallet.NewMockRepository(ctrl)
	mockTransactionRepo := transaction.NewMockRepository(ctrl)
	mockThirdPartyRepo := thirdparty.NewMockService(ctrl)

	// Create a sample request
	req := Request{
		Amount: 100.0,
		UserID: "user123",
		Type:   "debit",
	}

	// Create a sample user
	mockUser := user.User{
		ID: "user123",
	}

	// Create a sample wallet
	mockWallet := wallet.Wallet{
		UserID:  "user123",
		Balance: 20000, // $200.00 in cents
	}

	// Create a sample transaction
	mockTransaction := transaction.Transaction{
		UserID:    "user123",
		RequestID: uuid.NewString(),
		Reference: "ref123",
		Type:      "debit",
		Amount:    10000, // $100.00 in cents
	}

	// Set up expectations
	mockUserRepo.EXPECT().GetUserByID(req.UserID).Return(&mockUser, nil)
	mockWalletRepo.EXPECT().GetWalletByUserID(req.UserID).Return(&mockWallet, nil)
	mockTransactionRepo.EXPECT().Create(gomock.Any()).Return(&mockTransaction, nil)
	mockWalletRepo.EXPECT().DebitWallet(&mockWallet, gomock.Any(), mockTransaction.Amount).Return(&mockWallet, nil).Times(0)
	mockThirdPartyRepo.EXPECT().MakePayment(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
	mockTransactionRepo.EXPECT().UpdateTransactionToFailed(mockTransaction.ID).Return(&mockTransaction, nil)

	// Create the service with mocked dependencies
	svc := service{
		userRepo:          mockUserRepo,
		walletRepo:        mockWalletRepo,
		transactionRepo:   mockTransactionRepo,
		thirdPartyService: mockThirdPartyRepo,
	}

	// Call the method
	resp, err := svc.HandleTransactionRequest(req)

	// Check the result

	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Failed to make payment", resp.Message)

}
