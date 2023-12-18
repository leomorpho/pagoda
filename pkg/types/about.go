package types

import "html/template"

type (
	AboutData struct {
		ShowCacheWarning bool
		FrontendTabs     []AboutTab
		BackendTabs      []AboutTab
	}

	AboutTab struct {
		Title string
		Body  template.HTML
	}
)
