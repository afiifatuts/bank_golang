package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all function to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db), //its generate by SQLC and create and return queries object
	}
}

// execTx executes a function within a database transaction
// 1. add function to Store
// 2. to execute a generate database transaction
// 3. it takes a context
// 4. and callback function as input
// 5. return error
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx) //its same New function but instead passing sql.DB we pass db.BeginTx
	err = fn(q)
	if err != nil {
		// if rollback error return the message
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, rb err : %v ", err, rbErr)
		}
		//otherwise just return the error
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the result of the transfer transaction
type TransferTxResult struct {
	Transfer      Transfers `json:"tranfer"`
	FromAccountID Accounts  `json:"from_account_id"`
	ToAccountID   Accounts  `json:"to_account_id"`
	FromEntry     Entries   `json:"from_entry"`
	ToEntry       Entries   `json:"to_entry"`
}

// TransferTx perform a money transfer from one account to another
// 1. create transfer record
// 2. add account entries
// 3. update accounts balance whithin a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		//create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//create account entries - FromAccount
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}
		//create account entries - ToAccount
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		//Todo: update accounts balance 
		return nil
	})

	return result, err
}
