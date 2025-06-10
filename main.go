package main

import (
	"Syllybea/Render"
	"Syllybea/handler"
	"Syllybea/repository"
	"Syllybea/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"log"
	"os"
)

// TemplateRenderer is a custom renderer for Echo using the Go html/template package.
type TemplateRenderer struct {
	Templates *template.Template
}

// Render renders a template document.
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:admin@tcp(localhost:3306)/syllabus"
	}
	store, err := storage.NewStorage(dsn)
	if err != nil {
		log.Fatalf("Could not create storage: %v", err)
	}
	defer store.DB.Close()

	repo := repository.NewRepository(store.DB)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = Render.NewTemplate()

	e.Static("/static", "static")

	handler.RegisterRoutes(e, repo)

	e.Logger.Fatal(e.Start(":9090"))
}
