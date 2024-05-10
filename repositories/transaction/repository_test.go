package transaction

import (
	"p-system/tests"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionRepository(t *testing.T) {
	db := tests.StartDB(t)

	repo := NewRepository(db)

	//run seeds
	err := tests.Seed(db)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestCreateTransaction_Success", func(t *testing.T) {

		newTransaction := NewTransaction("d164e69d-26f5-448d-a18c-baeae517d9f2", "d164e69d-26f5-448d-a18c-baeae517d9f2", "newref", "credit", 1000)

		createdTransaction, err := repo.Create(newTransaction)

		require.NoError(t, err)

		require.Equal(t, newTransaction.UserID, createdTransaction.UserID)

	})
	t.Run("TestCreateTransaction_DuplicateReference", func(t *testing.T) {
		// Assuming you've seeded the database with a transaction having the reference "newref"
		duplicateTransaction := NewTransaction("d164e69d-26f5-448d-a18c-baeae517d9f2", "d164e69d-26f5-448d-a18c-baeae517d9f2", "newref", "credit", 1000)

		_, err := repo.Create(duplicateTransaction)

		require.Error(t, err)
		require.EqualError(t, err, "transaction already exists")
	})

	t.Run("TestGetTransactionByReference_NotFound", func(t *testing.T) {
		// Assuming "nonexistentref" does not exist in the seeded data
		_, err := repo.GetTransactionByReference("nonexistentref")

		require.Error(t, err)
		require.EqualError(t, err, "transaction not found")
	})

	t.Run("TestGetTransactionByReference_Success", func(t *testing.T) {
		// Assuming "newref" exists in the seeded data
		transaction, err := repo.GetTransactionByReference("newref")

		require.NoError(t, err)
		require.Equal(t, "newref", transaction.Reference)
	})

	t.Run("TestUpdateTransactionToFailed_Success", func(t *testing.T) {
		// Assuming "newref" exists in the seeded data
		transaction, err := repo.GetTransactionByReference("newref")

		require.NoError(t, err)
		require.Equal(t, "newref", transaction.Reference)

		updatedTransaction, err := repo.UpdateTransactionToFailed(transaction.ID)

		require.NoError(t, err)
		require.Equal(t, "failed", updatedTransaction.Status)
	})
}
