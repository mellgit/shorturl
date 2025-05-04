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

// ShortenHandler
// @Summary      ShortenHandler
// @Description  ShortenHandler create alias
// @Security ApiKeyAuth
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Param 		 request body ShortenRequest true "body"
// @Success      200 {object} ShortenResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/shorten [post]
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
		h.Logger.WithFields(log.Fields{
			"action": "CreateShortURL",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	response := ShortenResponse{
		ShortURL: ctx.BaseURL() + "/s/" + result.Alias,
		Expires:  result.ExpiresAt,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

// Protected
// @Summary      Protected
// @Description  Protected check authorized user
// @Security ApiKeyAuth
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/protected [get]
func (h *Handler) Protected(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(MessageResponse{Message: "authorized user"})
}

// Stats
// @Summary      Stats
// @Description  Stats for clicks on url
// @Security ApiKeyAuth
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Param        alias path string true "alias"
// @Success      200 {object} Count
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/stats/{alias} [get]
func (h *Handler) Stats(ctx *fiber.Ctx) error {

	alias := ctx.Params("alias")
	count, err := h.service.Stats(alias)
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "Stats",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}
	return ctx.Status(fiber.StatusOK).JSON(Count{Count: count})
}

// List
// @Summary      List
// @Description  List
// @Security ApiKeyAuth
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Success      200 {array} URL
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/shorten/list [get]
func (h *Handler) List(ctx *fiber.Ctx) error {

	urls, err := h.service.List()
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "List",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}
	return ctx.Status(fiber.StatusOK).JSON(urls)
}

// DeleteUrl
// @Summary      DeleteUrl
// @Description  DeleteUrl
// @Security ApiKeyAuth
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Param        alias path string true "alias"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/shorten/{alias} [delete]
func (h *Handler) DeleteUrl(ctx *fiber.Ctx) error {

	alias := ctx.Params("alias")
	if err := h.service.Delete(alias); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "DeleteUrl",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}
	return ctx.Status(fiber.StatusOK).JSON(MessageResponse{Message: "url deleted"})
}

// UpdateAlias
// @Summary      UpdateAlias
// @Description  UpdateAlias update only alias
// @Security ApiKeyAuth
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Param        alias path string true "alias"
// @Param 		 request body UpdateAliasRequest true "body"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/shorten/{alias} [patch]
func (h *Handler) UpdateAlias(ctx *fiber.Ctx) error {

	alias := ctx.Params("alias")
	payload := UpdateAliasRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ctx.BodyParser",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	if err := h.service.UpdateAlias(alias, payload.Alias); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "UpdateAlias",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}
	return ctx.Status(fiber.StatusOK).JSON(MessageResponse{Message: "url updated"})
}
