package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/panurujz/calculate-term-sheet/services"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is smooth.")
	})

	g := e.Group("/term-sheet")
	g.POST("/calculate", services.CalculateTs)

	e.Logger.Fatal(e.Start(":3001"))
}
