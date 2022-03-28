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

func Setup() *fiber.App {
	app := fiber.New()
	router.BlogRouter(app)

	return app
}

func main() {
	db := db.Connect()
	app := Setup()

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

	<-c
	fmt.Println("Gracefully shutting down...")
	app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// NOTE: Gorm mysql connection close
	mySqlDb, _ := db.DB()
	mySqlDb.Close()
	fmt.Println("Fiber was successful shutdown.")
}
