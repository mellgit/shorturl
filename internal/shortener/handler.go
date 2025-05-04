package shortener

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mellgit/shorturl/internal/middleware"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	service Service
	Logger  *log.Entry
}

func NewHandler(service Service, logger *log.Entry) *Handler {
	return &Handler{service, logger}
}

func (h *Handler) GroupHandler(app *fiber.App) {
	group := app.Group("/api", middleware.JWTProtected())
	group.Post("/shorten", h.ShortenHandler)
	group.Get("/shorten/list", h.List)
	group.Delete("/shorten/:alias", h.DeleteUrl)
	group.Patch("/shorten/:alias", h.UpdateAlias)
	group.Get("/protected", h.Protected)
	group.Get("/stats/:alias", h.Stats)
}

func (h *Handler) ShortenHandler(ctx *fiber.Ctx) error {
	payload := ShortenRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ctx.BodyParser",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	if payload.TTLHours == 0 {
		payload.TTLHours = 24
	}

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64)) // jwt numbers → float64

	result, err := h.service.CreateShortURL(userID, payload.URL, payload.Custom, payload.TTLHours)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"short_url":  ctx.BaseURL() + "/s/" + result.Alias,
		"expires_at": result.ExpiresAt,
	})
}

func (h *Handler) Protected(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"message": "authorized user"})
}

func (h *Handler) Stats(ctx *fiber.Ctx) error {

	alias := ctx.Params("alias")
	count, err := h.service.Stats(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(fiber.Map{
		"count": count,
	})

}

func (h *Handler) List(ctx *fiber.Ctx) error {

	urls, err := h.service.List()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(urls)
}

func (h *Handler) DeleteUrl(ctx *fiber.Ctx) error {

	alias := ctx.Params("alias")
	if err := h.service.Delete(alias); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "shortened url deleted",
	})
}

func (h *Handler) UpdateAlias(ctx *fiber.Ctx) error {

	alias := ctx.Params("alias")
	payload := UpdateAliasRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.service.UpdateAlias(alias, payload.Alias); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "shortened url updated",
	})
}
