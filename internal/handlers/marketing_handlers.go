package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

//nolint:wrapcheck //выводим результат, нет смысла оборачивать ошибки
func (h *MarketingHandler) UploadProducts(c *fiber.Ctx) error {
	ctx := c.Context()

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": fmt.Sprintf("failed to retrieve file: %s", err.Error())})
	}

	loadResult, err := h.service.SaveProducts(ctx, file)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": fmt.Sprintf("failed to process archive: %s", err.Error())})
	}

	return c.Status(http.StatusOK).JSON(loadResult)
}

//nolint:wrapcheck //выводим результат, нет смысла оборачивать ошибки
func (h *MarketingHandler) LoadProducts(c *fiber.Ctx) error {
	ctx := c.Context()

	zipData, err := h.service.LoadProducts(ctx)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": fmt.Sprintf("failed to laod products: %s", err.Error())})
	}

	// Устанавливаем заголовки для ответа
	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", "attachment; filename=data.zip")

	return c.Status(http.StatusOK).Send(zipData)
}
