package Render

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

// TemplateRenderer is a custom renderer for Echo that uses the Go html/template package.
type TemplateRenderer struct {
	Templates *template.Template
}

// Render implements echo's Renderer interface.
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *TemplateRenderer {
	funcs := template.FuncMap{
		"add1": func(i int) int { return i + 1 },
		"contains": func(arr []string, item string) bool {
			for _, v := range arr {
				if v == item {
					return true
				}
			}
			return false
		},
	}
	return &TemplateRenderer{
		Templates: template.Must(template.New("").Funcs(funcs).ParseGlob("views/*.html")),
	}
}
