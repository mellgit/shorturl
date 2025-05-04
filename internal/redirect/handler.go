package redirect

import (
	"github.com/gofiber/fiber/v2"
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
	group.Get("/s/:alias", h.RedirectHandler)
}

func (h *Handler) RedirectHandler(ctx *fiber.Ctx) error {

	alias := ctx.Params("alias")
	ip := ctx.IP()
	ua := ctx.Get("User-Agent")

	original, err := h.service.ResolveAndTrack(alias, ip, ua)
	if err != nil {
		h.Logger.WithFields(log.Fields{
			"action": "ResolveAndTrack",
		}).Errorf("%v", err)
		msgErr := ErrorResponse{Error: err.Error()}
		return ctx.Status(fiber.StatusBadRequest).JSON(msgErr)
	}

	//return ctx.Redirect(original, fiber.StatusFound)
	return ctx.Status(fiber.StatusOK).JSON(Original{Original: original})
}
