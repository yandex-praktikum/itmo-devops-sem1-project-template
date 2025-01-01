package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

//nolint:wrapcheck //выводим результат, нет смысла оборачивать ошибки
func (h *MarketingHandler) UploadProducts(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("failed to retrieve file: %s", err.Error()))
	}

	loadResult, err := h.service.SaveProducts(context.Background(), file)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("failed to process archive: %s", err.Error()))
	}

	return c.Status(http.StatusOK).JSON(loadResult)
}
