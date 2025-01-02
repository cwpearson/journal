package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cwpearson/journal/entries"
	"github.com/cwpearson/journal/ollama"
)

type CreateForm struct {
	UserText string `form:"userText"`
}

func CreatePost(c echo.Context) error {
	submission := new(CreateForm)

	// Bind form data to struct
	if err := c.Bind(submission); err != nil {
		return c.String(http.StatusBadRequest, "Invalid submission")
	}

	// Validate submission
	if submission.UserText == "" {
		return c.String(http.StatusBadRequest, "Text cannot be empty")
	}

	if _, err := entries.Create(ollama.NewClientFromConfig(), submission.UserText); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	return c.Redirect(http.StatusSeeOther, "/create")
}

func CreateGet(c echo.Context) error {
	return c.Render(http.StatusOK, "edit.html", map[string]interface{}{
		"title": "New Entry",
		"mode":  "create",
	})
}
