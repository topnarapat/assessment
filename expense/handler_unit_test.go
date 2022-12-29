package expense

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAllExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"}))
	mock.ExpectQuery("SELECT (.+) FROM expenses").WillReturnRows(rows)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := ExpenseHandler{db}

	expected := `[{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}]`

	if assert.NoError(t, h.GetExpensesHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetExpenseByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "americano", 89, "no sugar", pq.Array([]string{"beverage", "coffee"}))
	mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses WHERE id = (.+)").WillReturnRows(rows)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := ExpenseHandler{db}

	expected := `{"id":1,"title":"americano","amount":89,"note":"no sugar","tags":["beverage","coffee"]}`

	if assert.NoError(t, h.GetExpenseHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestCreateExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO expenses").WillReturnRows(rows)

	body := `{"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := ExpenseHandler{db}

	expected := `{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}`

	if assert.NoError(t, h.CreateExpenseHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestUpdateExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE expenses").
		ExpectExec().
		WithArgs(1, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"})).
		WillReturnResult(sqlmock.NewResult(0, 1))

	e := echo.New()
	body := `{"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := ExpenseHandler{db}

	expected := `{"id":1,"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]}`

	if assert.NoError(t, h.UpdateExpenseHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}
