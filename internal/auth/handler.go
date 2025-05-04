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
	group.Post("/register", h.Register)
}

func (h *Handler) Login(ctx *fiber.Ctx) error {

	payload := LoginRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ctx.BodyParser",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	token, err := h.service.Login(payload.Email, payload.Password)
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "Login",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	return ctx.Status(fiber.StatusOK).JSON(Token{Token: token})

}
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
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
