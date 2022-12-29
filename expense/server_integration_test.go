//go:build integration
// +build integration

package expense

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const serverPort = 80

func TestITGetExpenses(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewExpenseHandler(db)

		e.GET("/expenses", h.GetExpensesHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "[{\"id\":1,\"title\":\"buy a new phone\",\"amount\":39000,\"note\":\"buy a new phone\",\"tags\":[\"gadget\",\"shopping\"]}]"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestITGetExpense(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewExpenseHandler(db)

		e.GET("/expenses/:id", h.GetExpenseHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/1", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "{\"id\":1,\"title\":\"buy a new phone\",\"amount\":39000,\"note\":\"buy a new phone\",\"tags\":[\"gadget\",\"shopping\"]}"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestITCreateExpense(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewExpenseHandler(db)

		e.POST("/expenses", h.CreateExpenseHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := bytes.NewBufferString(`{
		"title": "buy a new phone",
		"amount": 39000,
		"note": "buy a new phone",
		"tags": ["gadget", "shopping"]
	}`)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), reqBody)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "{\"id\":2,\"title\":\"buy a new phone\",\"amount\":39000,\"note\":\"buy a new phone\",\"tags\":[\"gadget\",\"shopping\"]}"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestITUpdateExpense(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewExpenseHandler(db)

		e.PUT("/expenses/:id", h.UpdateExpenseHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := bytes.NewBufferString(`{
		"title": "apple smoothie",
		"amount": 89,
		"note": "apple smoothie",
		"tags": ["beverage"]
	}`)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/expenses/2", serverPort), reqBody)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "{\"id\":2,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"apple smoothie\",\"tags\":[\"beverage\"]}"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
