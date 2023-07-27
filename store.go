package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type store struct {
	db *sql.DB
}

// Store : Global store variable
var Store store

// initPostgresConnection : initializes Postgres connection
func initPostgresConnection() {
	var err error
	var db *sql.DB
	connStr := os.Getenv("POSTGRES_CONN_STR")
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		lg.Fatal().Msg("Error Opening connection to the database")
	}

	if err = db.Ping(); err != nil {
		lg.Fatal().Msg("Failed to ping database")
	}

	Store.db = db
	lg.Info().Msg("The database is connected")
}

func (s *store) saveTransaction(id int64, amount float64, transactionType string, parent *int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertStmt := `insert into "public"."transactions" ("id","amount", "type", "parent_id") values($1, $2, $3, $4)`
	_, err := s.db.ExecContext(ctx, insertStmt, id, amount, transactionType, parent)
	if err != nil {
		lg.Err(err).Int64("id", id).Float64("amount", amount).Str("transactionType", transactionType).Interface("parent", parent).Msg("Error inserting transaction to DB")
		return err
	}

	return err
}

func (s *store) getTransactionSum(ID int64) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var totalAmount float64
	var responseID string

	rows, err := s.db.QueryContext(ctx, `
	WITH RECURSIVE cte AS (
		SELECT id AS final_id, id, amount
		FROM transactions
		WHERE id = $1
		UNION ALL
		SELECT cte.final_id, tc.id, tc.amount
		FROM cte
		JOIN transactions tc ON tc.parent_id = cte.id
	  )
	  cycle id
		set is_cycle
 		using path
	  SELECT final_id, SUM(amount)
	  FROM cte where is_cycle=FALSE
	  GROUP BY final_id;
	`, ID)
	if err != nil {
		lg.Err(err).Int64("id", ID).Msg("Error getting SUM of transaction from DB")
		return totalAmount, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&responseID, &totalAmount)
	}
	return totalAmount, err
}

func (s *store) getTransactionByID(ID int64) (transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var t transaction

	rows, err := s.db.QueryContext(ctx, `SELECT "amount", "type", "parent_id" FROM  "public"."transactions" where id = $1 limit 1`, ID)
	if err != nil {
		lg.Err(err).Int64("id", ID).Msg("Error querying transaction from DB")
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&t.Amount, &t.Type, &t.Parent)
	}
	return t, err
}

func (s *store) getTransactionByType(t string) ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var txs []int64

	rows, err := s.db.QueryContext(ctx, `SELECT ID FROM "public"."transactions" where type = $1`, t)
	if err != nil {
		lg.Err(err).Str("type", t).Msg("Error querying transaction from DB")
		return txs, err
	}
	defer rows.Close()
	for rows.Next() {
		var t int64
		err = rows.Scan(&t)
		txs = append(txs, t)
	}
	return txs, err
}

func (s *store) deleteTransaction(ID int64, t string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertStmt := `delete from transactions where ID = $1 OR type = $2;`
	_, err := s.db.ExecContext(ctx, insertStmt, ID, t)
	if err != nil {
		lg.Err(err).Int64("id", ID).Msg("Error deleting transaction from DB")
		return err
	}

	return err
}

// Ping : Pings server for availability
func (s *store) Ping() error {
	return s.db.Ping()
}
