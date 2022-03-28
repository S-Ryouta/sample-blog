package api

import (
	"fmt"
	"github.com/S-Ryouta/sample-blog/db"
	"github.com/S-Ryouta/sample-blog/models"
	"github.com/S-Ryouta/sample-blog/serializers/entities"
	"github.com/gofiber/fiber/v2"
)

func GetEntries(c *fiber.Ctx) error {
	db := db.Connect()
	entries, err := models.SelectEntries(db)
	if err != nil {
		fmt.Println("failed to get entries", err)
	}

	return c.Status(fiber.StatusOK).JSON(entities.IndexSerializer(entries))
}

func GetEntry(c *fiber.Ctx) error {
	db := db.Connect()
	entry, err := models.FindEntry(db, c.Params("id"))
	if err != nil {
		fmt.Println("failed to get entry", err)
	}
	return c.Status(fiber.StatusOK).JSON(entry)
}
