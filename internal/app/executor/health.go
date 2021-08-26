package executor

import "github.com/gofiber/fiber/v2"

func (p TemplateExecutor) HealthEndpoint() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("OK 👋!")
	}
}
