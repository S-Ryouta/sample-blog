package main

import (
	"context"
	"fmt"
	"github.com/S-Ryouta/sample-blog/db"
	"github.com/S-Ryouta/sample-blog/job"
	"github.com/S-Ryouta/sample-blog/router"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db := db.Connect()
	app := fiber.New()
	router.BlogRouter(app)

	go func() {
		if err := app.Listen(":8000"); err != nil {
			log.Panic("Server closed with error: ", err)
		}
	}()

	// NOTE: contentful request
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	job.FetchContentful(ctx)

	// NOTE: Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// NOTE: Gorm mysql connection close
	mySqlDb, _ := db.DB()
	mySqlDb.Close()
	fmt.Println("Fiber was successful shutdown.")
}
