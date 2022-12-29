package expense

import "database/sql"

type DB interface {
	Query(query string, args ...any) (sql.Rows, error)
	QueryRow(query string, args ...any) sql.Row
	Exec(args ...any) (sql.Result, error)
}

type ExpenseHandler struct {
	db *sql.DB
}

func NewExpenseHandler(db *sql.DB) *ExpenseHandler {
	return &ExpenseHandler{db}
}
