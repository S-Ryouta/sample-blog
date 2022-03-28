package router

import (
	"github.com/S-Ryouta/sample-blog/controllers"
	"github.com/gofiber/fiber/v2"
)

func BlogRouter(app fiber.Router) {
	app.Get("/entries", func(c *fiber.Ctx) error {
		return controllers.GetEntries(c)
	})
	app.Get("/entries/:id", func(c *fiber.Ctx) error {
		return controllers.GetEntry(c)
	})
}
