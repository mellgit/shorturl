package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mellgit/shorturl/internal/middleware"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Handler struct {
	service Service
	Logger  *log.Entry
}

func NewHandler(service Service, logger *log.Entry) *Handler {
	return &Handler{service, logger}
}

func (h *Handler) GroupHandler(app *fiber.App) {
	group := app.Group("/api/users", middleware.JWTProtected())
	group.Get("", h.ListUsers)
	group.Get("/:id", h.GetUserByID)
	//group.Delete("/:id", h.DeleteUserByID)
	//group.Patch("/:id", h.UpdateUserByID)

}

func (h *Handler) ListUsers(ctx *fiber.Ctx) error {

	listUsers, err := h.service.ListUsers()
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "SomeUserService.GetAllUsers",
		}).Errorf("%v", err)
		msgErr := Error{Error: err.Error()}
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(msgErr)
	}
	return ctx.Status(fiber.StatusOK).JSON(listUsers)
}

func (h *Handler) GetUserByID(ctx *fiber.Ctx) error {

	id := ctx.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		msgErr := Error{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	user, err := h.service.GetUserByID(int64(idInt))
	if err != nil {
		msgErr := Error{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)

	}
	return ctx.Status(fiber.StatusOK).JSON(user)
}
