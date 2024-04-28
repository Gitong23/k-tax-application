package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gitong23/assessment-tax/config"
	"github.com/Gitong23/assessment-tax/postgres"
	"github.com/Gitong23/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	config := config.New()
	p, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Validator = tax.NewValidator()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	handler := tax.NewHandler(p)
	e.POST("/tax/calculations", handler.Tax)
	e.POST("/tax/calculations/upload-csv", handler.UploadCsv)

	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == config.Credentials.Username && password == config.Credentials.Password {
			return true, nil
		}
		return false, c.JSON(http.StatusUnauthorized, tax.Err{Message: "Unauthorized"})
	}))

	g.POST("/deductions/personal", handler.UpdateInitPersonalDeduct)
	g.POST("/deductions/k-receipt", handler.UpdateMaxKreceiptDeduct)

	// Graceful shutdown
	go func() {
		port := fmt.Sprintf(":%s", config.Server.Port)
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Info("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
