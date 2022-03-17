package router

import (
	"github.com/S-Ryouta/sample-blog/api"
	"github.com/gofiber/fiber/v2"
)

func BlogRouter(app fiber.Router) {
	app.Get("/entities", func(c *fiber.Ctx) error {
		return api.GetEntries(c)
	})
	app.Get("/entities/:id", func(c *fiber.Ctx) error {
		return api.GetEntry(c)
	})
}
