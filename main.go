package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/noobbrother112/run_monitor/db"
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

func handleHome(c echo.Context) error {
	db.Setdb()
	tData := data{
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
	return c.Render(http.StatusOK, "home.html", tData)
}

func addLog(c echo.Context) error {
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {

		log.Error("empty json body")
		return nil
	}
	// 들어온 json data 검증
	if jsonBody["code"] != nil {
		db.AddLog(jsonBody)
	}
	return c.String(200, "done")
}

func main() {
	db.Setdb()
	t := &Template{
		templates: template.Must(template.ParseGlob("html/*.html")),
	}

	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/log", addLog)
	e.Renderer = t
	e.Logger.Fatal(e.Start(":1323")) // localhost:1323
}
