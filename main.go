package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/cwpearson/journal/config"
	"github.com/cwpearson/journal/database"
	"github.com/cwpearson/journal/handlers"
	"github.com/cwpearson/journal/models"
	"github.com/cwpearson/journal/ollama"
)

// Template renderer
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	// create config dir
	err := os.MkdirAll(config.ConfigDir(), 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	database.Init()
	database.Get().AutoMigrate(&models.Entry{}, &models.Tag{}, &ollama.Record{})
	handlers.Init()

	client := ollama.NewClientFromConfig()
	if err := client.Pull(); err != nil {
		log.Println(err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Templates
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	e.GET("/login", handlers.LoginGet)
	e.POST("/login", handlers.LoginPost)

	e.GET("/", handlers.RootGet, handlers.AuthMiddleware)
	e.GET("/history", handlers.HistoryGet, handlers.AuthMiddleware)
	e.GET("/:year/:month/:day", handlers.EditGet, handlers.AuthMiddleware)
	e.POST("/:year/:month/:day", handlers.EditPost, handlers.AuthMiddleware)
	e.POST("/delete/:id", handlers.DeletePost, handlers.AuthMiddleware)
	e.POST("/logout", handlers.LogoutPost, handlers.AuthMiddleware)

	staticGroup := e.Group("/static")
	staticGroup.Static("/", "static")

	addr := fmt.Sprintf(":%s", config.Port())
	e.Logger.Fatal(e.Start(addr))
}
