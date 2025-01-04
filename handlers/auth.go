package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cwpearson/journal/config"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := store.Get(c.Request(), "session")
		if err != nil {
			log.Println("AuthMiddleware: unable to retrieve session. Redirect to /login")
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		_, ok := session.Values["logged_in"]
		if !ok {
			log.Println("AuthMiddleware: session does not contain logged_in. Redirect to /login")
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		return next(c)
	}
}

func LoginPost(c echo.Context) error {
	password := c.FormValue("password")

	if password != config.Password() {
		return c.String(http.StatusUnauthorized, "Invalid password")
	}

	session, err := store.Get(c.Request(), "session")
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "LoginPost: Unable to retrieve session")
	}

	session.Values["logged_in"] = true
	err = session.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to save session")
	}

	session, _ = store.Get(c.Request(), "session")
	_, ok := session.Values["logged_in"]
	if !ok {
		return c.String(http.StatusInternalServerError, "logged_in was not saved as expected")
	}

	fmt.Println("loginPostHandler: redirect to /")
	return c.Redirect(http.StatusSeeOther, "/")
}

func LoginGet(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func LogoutPost(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session")
	delete(session.Values, "logged_in")
	session.Save(c.Request(), c.Response().Writer)
	return c.Redirect(http.StatusSeeOther, "/")
}
