package handlers

import (
	"log"

	"github.com/cwpearson/journal/config"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func Init() error {
	log.Println("handlers.Init()...")
	key := config.SessionKey()
	store = sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7, // seconds
		HttpOnly: true,
		Secure:   config.SessionSecure(),
	}

	return nil
}
