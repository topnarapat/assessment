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

func main() {
	fmt.Println("Please use server.go for main file")

	expense.InitDB()

	e := echo.New()
	e.GET("/health", healthHandler)
	e.POST("/expenses", expense.CreateExpenseHandler)

	fmt.Println("start at port:", os.Getenv("PORT"))
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
