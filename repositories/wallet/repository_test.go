package wallet

import (
	"p-system/repositories/transaction"
	"p-system/tests"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWalletRepository(t *testing.T) {
	db := tests.StartDB(t)

	repo := NewRepository(db)

	// Run seeds
	err := tests.Seed(db)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestCreateWallet_Success", func(t *testing.T) {
		newWallet := &Wallet{
			UserID:    "d164e69d-26f5-448d-a18c-baeae517d9f2",
			Balance:   1000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		createdWallet, err := repo.Create(newWallet)

		require.NoError(t, err)
		require.Equal(t, newWallet.UserID, createdWallet.UserID)
	})

	t.Run("TestGetWalletByUserID_Success", func(t *testing.T) {
		expectedUserID := "d164e69d-26f5-448d-a18c-baeae517d9f2"

		wallet, err := repo.GetWalletByUserID(expectedUserID)

		require.NoError(t, err)
		require.NotNil(t, wallet)
		require.Equal(t, expectedUserID, wallet.UserID)
	})

	t.Run("TestCreditWallet_Success", func(t *testing.T) {
		wallet := &Wallet{
			ID:        "d164e69d-26f5-448d-a18c-baeae517d991",
			UserID:    "d164e69d-26f5-448d-a18c-baeae517d9f2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		transaction := transaction.Transaction{
			ID:     "d164e69d-26f5-448d-a18c-baeae517d9f5",
			Amount: 500,
		}

		updatedWallet, err := repo.CreditWallet(wallet, transaction, 5000)

		require.NoError(t, err)
		require.NotNil(t, updatedWallet)
		require.Equal(t, int64(5000), updatedWallet.Balance)
	})

	t.Run("TestDebitWallet_Success", func(t *testing.T) {
		wallet := &Wallet{
			ID:        "d164e69d-26f5-448d-a18c-baeae517d991",
			UserID:    "d164e69d-26f5-448d-a18c-baeae517d9f2",
			Balance:   1000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		transaction := transaction.Transaction{
			ID: "d164e69d-26f5-448d-a18c-baeae517d9f5",
		}

		updatedWallet, err := repo.DebitWallet(wallet, transaction, 500)

		require.NoError(t, err)
		require.NotNil(t, updatedWallet)
		require.Equal(t, int64(4500), updatedWallet.Balance)
	})
}
