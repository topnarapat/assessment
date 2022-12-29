package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/topnarapat/assessment/expense"
)

func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func checkAuthHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
		}
		if authHeader != "November 10, 2009" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is invalid")
		}
		return next(c)
	}
}

func main() {
	fmt.Println("Please use server.go for main file")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	e := echo.New()

	e.Use(checkAuthHeader)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler := expense.NewExpenseHandler(db)

	handler.InitDB()

	e.GET("/health", healthHandler)
	e.POST("/expenses", handler.CreateExpenseHandler)
	e.GET("/expenses/:id", handler.GetExpenseHandler)
	e.PUT("/expenses/:id", handler.UpdateExpenseHandler)
	e.GET("/expenses", handler.GetExpensesHandler)

	fmt.Println("start at port:", os.Getenv("PORT"))

	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
