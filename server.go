package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

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

	expense.InitDB()

	e := echo.New()

	e.Use(checkAuthHeader)

	e.GET("/health", healthHandler)
	e.POST("/expenses", expense.CreateExpenseHandler)
	e.GET("/expenses/:id", expense.GetExpenseHandler)
	e.PUT("/expenses/:id", expense.UpdateExpenseHandler)

	fmt.Println("start at port:", os.Getenv("PORT"))
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
