package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gitong23/assessment-tax/postgres"
	"github.com/Gitong23/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func main() {

	p, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	// e.Logger.SetLevel(log.INFO)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	taxHandler := tax.NewHandler(p)
	e.POST("/tax/calculation", taxHandler.CalTax)

	// Graceful shutdown
	go func() {
		port := fmt.Sprintf(":%s", os.Getenv("PORT"))
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
