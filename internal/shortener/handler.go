package shortener

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mellgit/shorturl/internal/middleware"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	service *Service
	Logger  *log.Entry
}

func NewHandler(service *Service, logger *log.Entry) *Handler {
	return &Handler{service, logger}
}

func (h *Handler) GroupHandler(app *fiber.App) {
	group := app.Group("/api", middleware.JWTProtected())
	group.Post("/shorten", h.ShortenHandler)
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
