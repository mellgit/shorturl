package redirect

import (
	"github.com/gofiber/fiber/v2"
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
	group.Get("/s/:alias", h.RedirectHandler)

}

func (h *Handler) RedirectHandler(ctx *fiber.Ctx) error {

	alias := ctx.Params("alias")
	ip := ctx.IP()
	ua := ctx.Get("User-Agent")
	h.Logger.Info(alias, ip, ua)

	original, err := h.service.ResolveAndTrack(alias, ip, ua)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	//return ctx.Redirect(original, fiber.StatusFound)
	return ctx.JSON(fiber.Map{"original": original})
}

//
//func RedirectHandler(s *Service) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		alias := c.Params("alias")
//		ip := c.IP()
//		ua := c.Get("User-Agent")
//
//		original, err := s.ResolveAndTrack(alias, ip, ua)
//		if err != nil {
//			return fiber.NewError(fiber.StatusNotFound, err.Error())
//		}
//
//		return c.Redirect(original, fiber.StatusFound)
//	}
//}
