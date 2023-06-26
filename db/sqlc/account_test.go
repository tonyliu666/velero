package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tonyliu666/simplebank/db/util"
)

// func TestCreateAccount(t *testing.T)
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)      //this command will check the err is nil or not
	require.NotEmpty(t, account) //account should not be empty object
	// additionally we should make a comparision between the actual object and defined object
	//defined object which is the same as arg(above)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt) //assure that timestamp is not zero

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	//in order to avoid code duplication,we need to use the CreateAccount function to add a new entry
	account1 := createRandomAccount(t)
	// get the entry ID
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	//check the error is nil
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	// check the account.CreatedAt time delta is smaller than one second
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	//everything should remain the same except for the balance
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	// check the account.CreatedAt time delta is smaller than one second
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	//it might occur errors because it is deleted.
	require.Error(t, err)
	//check what type of this error is
	require.EqualError(t, err, sql.ErrNoRows.Error())
	//and account2 should be empty
	require.Empty(t, account2)
}
func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	//means we only need the tail 5 elements from the index 5
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
