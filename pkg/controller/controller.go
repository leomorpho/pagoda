package controller

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/htmx"
	"github.com/mikestefanello/pagoda/pkg/middleware"
	"github.com/mikestefanello/pagoda/pkg/services"

	"github.com/labstack/echo/v4"
)

// Controller provides base functionality and dependencies to routes.
// The proposed pattern is to embed a Controller in each individual route struct and to use
// the router to inject the container so your routes have access to the services within the container
type Controller struct {
	// Container stores a services container which contains dependencies
	Container *services.Container
}

// NewController creates a new Controller
func NewController(c *services.Container) Controller {
	return Controller{
		Container: c,
	}
}

// RenderPage renders a Page as an HTTP response
func (c *Controller) RenderPage(ctx echo.Context, page Page) error {
	buf := &bytes.Buffer{}
	var err error

	// Page name is required
	if page.Name == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "page render failed due to missing name")
	}

	// Page component required
	if page.Component == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "page render failed due to missing component")
	}

	// Use the app name in configuration if a value was not set
	if page.AppName == "" {
		page.AppName = c.Container.Config.App.Name
	}

	// Check if this is an HTMX non-boosted request which indicates that only partial
	// content should be rendered
	if page.HTMX.Request.Enabled && !page.HTMX.Request.Boosted {
		// Render the templates only for the content portion of the page
		page.Component.Render(ctx.Request().Context(), buf)
	} else {
		// Render the templates for the Page
		// If the page Layout is set, that will be used to wrap the page component
		component := page.Component
		if page.Layout != nil {
			component = page.Layout(component, &page)
		}
		component.Render(ctx.Request().Context(), buf)
	}

	if err != nil {
		return c.Fail(err, "failed to parse and execute templates")
	}

	// Set the status code
	ctx.Response().Status = page.StatusCode

	// Set any headers
	for k, v := range page.Headers {
		ctx.Response().Header().Set(k, v)
	}

	// Apply the HTMX response, if one
	if page.HTMX.Response != nil {
		page.HTMX.Response.Apply(ctx)
	}

	// Cache this page, if caching was enabled
	c.cachePage(ctx, page, buf)

	return ctx.HTMLBlob(ctx.Response().Status, buf.Bytes())
}

// cachePage caches the HTML for a given Page if the Page has caching enabled
func (c *Controller) cachePage(ctx echo.Context, page Page, html *bytes.Buffer) {
	if !page.Cache.Enabled || page.IsAuth {
		return
	}

	// If no expiration time was provided, default to the configuration value
	if page.Cache.Expiration == 0 {
		page.Cache.Expiration = c.Container.Config.Cache.Expiration.Page
	}

	// Extract the headers
	headers := make(map[string]string)
	for k, v := range ctx.Response().Header() {
		headers[k] = v[0]
	}

	// The request URL is used as the cache key so the middleware can serve the
	// cached page on matching requests
	key := ctx.Request().URL.String()
	cp := middleware.CachedPage{
		URL:        key,
		HTML:       html.Bytes(),
		Headers:    headers,
		StatusCode: ctx.Response().Status,
	}

	err := c.Container.Cache.
		Set().
		Group(middleware.CachedPageGroup).
		Key(key).
		Tags(page.Cache.Tags...).
		Expiration(page.Cache.Expiration).
		Data(cp).
		Save(ctx.Request().Context())

	switch {
	case err == nil:
		ctx.Logger().Info("cached page")
	case !context.IsCanceledError(err):
		ctx.Logger().Errorf("failed to cache page: %v", err)
	}
}

// Redirect redirects to a given route name with optional route parameters
func (c *Controller) Redirect(ctx echo.Context, route string, routeParams ...any) error {
	url := ctx.Echo().Reverse(route, routeParams...)

	if htmx.GetRequest(ctx).Boosted {
		htmx.Response{
			Redirect: url,
		}.Apply(ctx)

		return nil
	} else {
		return ctx.Redirect(http.StatusFound, url)
	}
}

// Fail is a helper to fail a request by returning a 500 error and logging the error
func (c *Controller) Fail(err error, log string) error {
	return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("%s: %v", log, err))
}
