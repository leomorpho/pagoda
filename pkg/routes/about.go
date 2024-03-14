package routes

import (
	"html/template"

	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/pkg/types"
	"github.com/mikestefanello/pagoda/templates"
	"github.com/mikestefanello/pagoda/templates/layouts"
	"github.com/mikestefanello/pagoda/templates/pages"

	"github.com/labstack/echo/v4"
)

type (
	about struct {
		controller.Controller
	}
)

func (c *about) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = layouts.Main
	page.Name = templates.PageAbout
	page.Title = "About"
	page.Component = pages.About(&page)
	page.HTMX.Request.Boosted = true

	// This page will be cached!
	// page.Cache.Enabled = true
	page.Cache.Tags = []string{"page_about", "page:list"}

	// A simple example of how the Data field can contain anything you want to send to the templates
	// even though you wouldn't normally send markup like this
	page.Data = types.AboutData{
		ShowCacheWarning: true,
		FrontendTabs: []types.AboutTab{
			{
				Title: "HTMX",
				Body: template.HTML(
					`Completes HTML as a hypertext by providing attributes to AJAXify anything and much more. ` +
						`Visit <a class="text-blue-400" href="https://htmx.org/">htmx.org</a> to learn more.`,
				),
			},
			{
				Title: "Alpine.js",
				Body: template.HTML(
					`Drop-in, Vue-like functionality written directly in your markup. Visit ` +
						`<a class="text-blue-400" href="https://alpinejs.dev/">alpinejs.dev</a> to learn more.`,
				),
			},
			{
				Title: "TailwindCSS",
				Body: template.HTML(
					`Ready-to-use frontend components and styling that you can easily combine to build responsive web interface. ` +
						`Visit <a class="text-blue-400" href="https://tailwindcss.com/">tailwind.com</a> to learn more. ` +
						`Note that the Daisyui is also for easy styling through their premade color names, ` +
						`see <a class="text-blue-400" href="https://daisyui.com/docs/colors/#-2">daisyui.com</a> ` +
						`for more info.`,
				),
			},
			{
				Title: "Javascript",
				Body: template.HTML(
					`<p>A full vanilla JS build system that will run on HTMX swaps and allow you ` +
						`to create reusable bits of logic that can be housed independently of ` +
						`Templ templates.</p><div id="js-quiz-container" class="my-4"></div>`,
				),
			},
			{
				Title: "Svelte",
				Body: template.HTML(
					`<p>Create islands of high interactivity with a JS framework. Note that Svelte ` +
						`could be swapped for just about any JS framework you like.</p>` +
						`<div id="test-svelte-todo-list" class="my-4"></div>` +
						`<div id="test-multi-select"></div>`,
				),
			},
		},
		BackendTabs: []types.AboutTab{
			{
				Title: "Echo",
				Body: template.HTML(
					`High performance, extensible, minimalist Go web framework. ` +
						`Visit <a class="text-blue-400" href="https://echo.labstack.com/">echo.labstack.com</a> to learn more.`,
				),
			},
			{
				Title: "Ent",
				Body: template.HTML(
					`Simple, yet powerful ORM for modeling and querying data. Visit ` +
						`<a class="text-blue-400" href="https://entgo.io/">entgo.io</a> to learn more.`,
				),
			},
			{
				Title: "Templ",
				Body: template.HTML(
					`A language for writing HTML user interfaces in Go. Visit ` +
						`<a class="text-blue-400" href="https://templ.guide/">templ.guide</a> to learn more.`,
				),
			},
		},
	}

	return c.RenderPage(ctx, page)
}
