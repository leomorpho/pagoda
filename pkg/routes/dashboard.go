package routes

import (
	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"

	"github.com/labstack/echo/v4"
)

type (
	dashboard struct {
		controller.Controller
	}
)

func (c *dashboard) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = layouts.Main
	page.Name = templates.PageDashboard
	page.Component = pages.Dashboard(&page)
	page.HTMX.Request.Boosted = true

	return c.RenderPage(ctx, page)
}
