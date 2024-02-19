// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"strings"

	"github.com/mikestefanello/pagoda/pkg/controller"
	"github.com/mikestefanello/pagoda/templates/helpers"
)

func Metatags(page *controller.Page) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(page.AppName)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/core.templ`, Line: 12, Col: 16}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(page.Title) > 0 {
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("| %s", page.Title))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/core.templ`, Line: 14, Col: 36}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><link rel=\"icon\" href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(helpers.File("favicon.png")))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(page.Metatags.Description) > 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<meta name=\"description\" content=\"{page.Metatags.Description}\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if len(page.Metatags.Keywords) > 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<meta name=\"keywords\" content=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(strings.Join(page.Metatags.Keywords, ", ")))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func CSS() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<link rel=\"stylesheet\" href=\"/files/svelte_bundle.css\"><link rel=\"stylesheet\" href=\"/files/styles_bundle.css\"><link rel=\"manifest\" href=\"/files/pwa/manifest.json\"><link href=\"https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.2.1/flowbite.min.css\" rel=\"stylesheet\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func JS() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var5 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var5 == nil {
			templ_7745c5c3_Var5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script src=\"https://unpkg.com/htmx.org@1.9.10\" integrity=\"sha384-D1Kt99CQMDuVetoL1lrYwg5t+9QdHe7NLX/SoJYkXDFfX37iInKRy5xLSi8nO7UC\" crossorigin=\"anonymous\"></script><script defer src=\"https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js\"></script><script src=\"/files/vanilla_bundle.js\"></script><script src=\"/files/svelte_bundle.js\"></script><script src=\"/service-worker.js\"></script><script src=\"https://unpkg.com/htmx.org/dist/ext/debug.js\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = darkModeSwitcher().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func darkModeSwitcher() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_darkModeSwitcher_5a44`,
		Function: `function __templ_darkModeSwitcher_5a44(){// On page load or when changing themes, best to add inline in ` + "`" + `head` + "`" + ` to avoid FOUC
    if (localStorage.getItem('color-theme') === 'darkmode' || (!('color-theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
        document.documentElement.classList.add('dark'); // For default tailwind dark: utility
		document.documentElement.setAttribute('data-theme', 'darkmode'); // For daisyui theme
		// Set the hover brightness dynamically
		document.documentElement.style.setProperty('--brightness-hover', 'var(--brightness-hover-dark)');
    } else {
        document.documentElement.classList.remove('dark') // For default tailwind dark: utility
		document.documentElement.setAttribute('data-theme', 'lightmode'); // For daisyui theme
		// Set the hover brightness dynamically
		document.documentElement.style.setProperty('--brightness-hover', 'var(--brightness-hover-light)');
    }
}`,
		Call:       templ.SafeScript(`__templ_darkModeSwitcher_5a44`),
		CallInline: templ.SafeScriptInline(`__templ_darkModeSwitcher_5a44`),
	}
}

func csrfJS(token string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_csrfJS_e551`,
		Function: `function __templ_csrfJS_e551(token){document.body.addEventListener('htmx:configRequest', function(evt)  {
        if (evt.detail.verb !== "get") {
            evt.detail.parameters['csrf'] = token;
        }
    })
}`,
		Call:       templ.SafeScript(`__templ_csrfJS_e551`, token),
		CallInline: templ.SafeScriptInline(`__templ_csrfJS_e551`, token),
	}
}

func htmxBeforeSwap() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_htmxBeforeSwap_7d49`,
		Function: `function __templ_htmxBeforeSwap_7d49(){document.body.addEventListener('htmx:beforeSwap', function(evt) {
        if (evt.detail.xhr.status >= 400){
            evt.detail.shouldSwap = true;
            evt.detail.target = htmx.find("body");
        }
    });
}`,
		Call:       templ.SafeScript(`__templ_htmxBeforeSwap_7d49`),
		CallInline: templ.SafeScriptInline(`__templ_htmxBeforeSwap_7d49`),
	}
}

func htmxOnLoad() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_htmxOnLoad_0da8`,
		Function: `function __templ_htmxOnLoad_0da8(){htmx.onLoad(function(content) {
			initializeJS();
			initializeAppSvelte();
	});
}`,
		Call:       templ.SafeScript(`__templ_htmxOnLoad_0da8`),
		CallInline: templ.SafeScriptInline(`__templ_htmxOnLoad_0da8`),
	}
}

func Footer(page *controller.Page) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var6 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var6 == nil {
			templ_7745c5c3_Var6 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if len(page.CSRF) > 0 {
			templ_7745c5c3_Err = csrfJS(page.CSRF).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = htmxBeforeSwap().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = htmxOnLoad().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func BeforeBodyEnd() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var7 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var7 == nil {
			templ_7745c5c3_Var7 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script src=\"https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.2.1/flowbite.min.js\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
