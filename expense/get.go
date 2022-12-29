package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *ExpenseHandler) GetExpenseHandler(c echo.Context) error {
	id := c.Param("id")
	row := h.db.QueryRow("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1", id)
	e := Expense{}
	err := row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}
}

func (h *ExpenseHandler) GetExpensesHandler(c echo.Context) error {
	rows, err := h.db.Query("SELECT * FROM expenses")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "can't query all expenses:" + err.Error()})
	}
	expenses := []Expense{}
	for rows.Next() {
		var e Expense
		err = rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: "can't scan expenses:" + err.Error()})
		}
		expenses = append(expenses, e)
	}

	return c.JSON(http.StatusOK, expenses)
}
