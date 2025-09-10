package rest

import "github.com/gofiber/fiber/v2"

func ErrorMessage(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(err.Error())
}

func InternalError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(err.Error())
}

func SuccessMessage(ctx *fiber.Ctx, message string, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": message,
		"data":    data,
	})
}
