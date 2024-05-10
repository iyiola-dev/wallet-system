package user

import (
	"p-system/tests"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUserByID(t *testing.T) {
	db := tests.StartDB(t)

	repo := NewRepository(db)

	//run seeds
	err := tests.Seed(db)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestGetUserByID_Success", func(t *testing.T) {
		expectedUser := &User{
			ID:    "d164e69d-26f5-448d-a18c-baeae517d9f2",
			Email: "john@ample.com",
		}

		user, err := repo.GetUserByID(expectedUser.ID)

		require.NoError(t, err)
		require.Equal(t, expectedUser.ID, user.ID)
		require.Equal(t, expectedUser.Email, user.Email)
	})

	t.Run("TestGetUserByID_NotFound", func(t *testing.T) {
		// Assuming "nonexistentid" does not exist in the seeded data
		_, err := repo.GetUserByID("d164e69d-26f5-448d-a18c-baeae517d912")

		require.Error(t, err)
		require.EqualError(t, err, "user not found")
	})
}
