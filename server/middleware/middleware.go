package middleware

import (
	"net/http"

	"github.com/justinas/alice"
)

func NewHandler() http.Handler {
	return alice.New(recoverHandler, authHandler).ThenFunc(logicHandler)
}

func recoverHandler() http.Handler {

}

func authHandler() http.Handler {

}

func logicHandler() {

}
