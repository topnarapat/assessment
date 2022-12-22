package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func UpdateExpenseHandler(c echo.Context) error {
	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	e := Expense{}
	// find by id
	err = c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := db.Prepare(`
	UPDATE expenses
	SET title=$2, amount=$3, note=$4, tags=$5
	WHERE id=$1
	`)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if _, err := stmt.Exec(rowID, e.Title, e.Amount, e.Note, pq.Array(e.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	e.ID = rowID
	return c.JSON(http.StatusOK, e)
}
