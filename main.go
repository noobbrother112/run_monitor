package main

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
	"net/http"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type data struct {
	Datafields map[string]interface{}
}

func home(c echo.Context) error {
	t_data := data{
		Datafields: map[string]any{
			"081": map[string]any{
				"id":   1,
				"char": 1,
			},
			"082": map[string]any{
				"id":   2,
				"char": 2,
			},
			"083": map[string]any{
				"id":   3,
				"char": 3,
			},
		},
	}
	return c.Render(http.StatusOK, "home.html", t_data)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("html/*.html")),
	}

	e := echo.New()
	e.GET("/", home).Name = "sukho"
	e.Renderer = t
	e.Logger.Fatal(e.Start(":1323")) // localhost:1323
}
