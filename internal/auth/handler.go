package auth

import (
	"github.com/gofiber/fiber/v2"
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
	group := app.Group("/auth")
	group.Post("/login", h.Login)
	group.Post("/logout", h.Logout)
	group.Post("/register", h.Register)
	group.Post("/refresh", h.RefreshToken)
}

// RefreshToken
// @Summary      RefreshToken
// @Description  RefreshToken
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Param 		 request body RefreshTokenRequest true "body"
// @Success      200 {object} AccessTokenResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/refresh [post]
func (h *Handler) RefreshToken(ctx *fiber.Ctx) error {

	payload := RefreshTokenRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ctx.BodyParser",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	accessToken, err := h.service.RefreshToken(payload.RefreshToken)
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "RefreshToken",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}

	return ctx.Status(fiber.StatusOK).JSON(accessToken)
}

// Login
// @Summary      Login
// @Description  Login user
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Param 		 request body LoginRequest true "body"
// @Success      200 {object} TokensResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/login [post]
func (h *Handler) Login(ctx *fiber.Ctx) error {

	payload := LoginRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ctx.BodyParser",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	tokens, err := h.service.Login(payload.Email, payload.Password)
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "Login",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}

	return ctx.Status(fiber.StatusOK).JSON(tokens)
}

// Register
// @Summary      Register
// @Description  Register new user
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Param 		 request body RegisterRequest true "body"
// @Success      204 {object} int
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/register/ [post]
func (h *Handler) Register(ctx *fiber.Ctx) error {

	payload := RegisterRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ctx.BodyParser",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	if err := h.service.Register(payload.Email, payload.Password); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "Register",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

// Logout
// @Summary      Logout
// @Description  Logout user
// @Tags         ShortUrl
// @Accept       json
// @Produce      json
// @Param 		 request body LogoutRequest true "body"
// @Success      200 {object} int
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/logout [post]
func (h *Handler) Logout(ctx *fiber.Ctx) error {

	payload := LogoutRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ctx.BodyParser",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	if err := h.service.Logout(payload.RefreshToken); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "Logout",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusInternalServerError).JSON(msgErr)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
