package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ahmed-deftoner/csrf-go/db"
	myjwt "github.com/ahmed-deftoner/csrf-go/server/middleware/myJwt"
	"github.com/ahmed-deftoner/csrf-go/server/templates"
	"github.com/justinas/alice"
)

func NewHandler() http.Handler {
	return alice.New(recoverHandler, authHandler).ThenFunc(logicHandler)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Panic("Recovered! Panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func authHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/restricted", "/logout", "/deleteuser":
			log.Println("In auth restricted section")

			// read cookies
			AuthCookie, authErr := r.Cookie("AuthToken")
			if authErr == http.ErrNoCookie {
				log.Println("Unauthorized attempt! No auth cookie")
				nullifyCookies(&w, r)
				// http.Redirect(w, r, "/login", 302)
				http.Error(w, http.StatusText(401), 401)
				return
			} else if authErr != nil {
				log.Panic("panic: %+v", authErr)
				nullifyCookies(&w, r)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			RefreshCookie, refreshErr := r.Cookie("RefreshToken")
			if refreshErr == http.ErrNoCookie {
				log.Println("Unauthorized attempt! No refresh cookie")
				nullifyCookies(&w, r)
				http.Redirect(w, r, "/login", 302)
				return
			} else if refreshErr != nil {
				log.Panic("panic: %+v", refreshErr)
				nullifyCookies(&w, r)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			requestCsrfToken := getCsrfToken(r)
			log.Println(requestCsrfToken)

			authTokenString, refreshTokenString, csrfSecret, err := myjwt.CheckAndRefreshTokens(AuthCookie.Value, RefreshCookie.Value, requestCsrfToken)
		default:
		}
	}
	return http.HandlerFunc(fn)
}

func logicHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/restricted":
		csrfSecret := getCsrfToken(r)
		templates.RenderTemplate(w, "restricted", &templates.Restricted{csrfSecret, "Hello Ahmed!"})

	case "/login":
		switch r.Method {
		case "/GET":
		case "/POST":
		default:
		}
	case "/register":
		switch r.Method {
		case "/GET":
			templates.RenderTemplate(w, "register", &templates.Register{false, ""})
		case "/POST":
			r.ParseForm()
			log.Println(r.Form)

			// check to see if the username is already taken
			_, uuid, err := db.FetchUserByUsername(strings.Join(r.Form["username"], ""))
			if err == nil {
				// templates.RenderTemplate(w, "register", &templates.RegisterPage{ true, "Username not available!" })
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				role := "user"
				uuid, err = db.StoreUser(strings.Join(r.Form["username"], ""), strings.Join(r.Form["password"], ""), role)
				if err != nil {
					http.Error(w, http.StatusText(500), 500)
				}
				log.Println("uuid: " + uuid)

				authTokenString, refreshTokenString, csrfSecret, err := myjwt.CreateNewTokens(uuid, role)
				if err != nil {
					http.Error(w, http.StatusText(500), 500)
				}

				setAuthandRefreshCookies(&w, authTokenString, refreshTokenString)
				w.Header().Set("X-CSRF-Token", csrfSecret)

				w.WriteHeader(http.StatusOK)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/logout":
		nullifyCookies(&w, r)
		http.Redirect(w, r, "/login", 302)
	case "/deleteuser":
	default:
	}
}

func nullifyCookies(w *http.ResponseWriter, r *http.Request) {
	authCookie := http.Cookie{
		Name:     "AuthToken",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(*w, &authCookie)

	refreshCookie := http.Cookie{
		Name:     "RefreshToken",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(*w, &refreshCookie)

	// if present, revoke the refresh cookie from our db
	RefreshCookie, refreshErr := r.Cookie("RefreshToken")
	if refreshErr == http.ErrNoCookie {
		// do nothing, there is no refresh cookie present
		return
	} else if refreshErr != nil {
		log.Panic("panic: %+v", refreshErr)
		http.Error(*w, http.StatusText(500), 500)
	}
	myjwt.RevokeRefreshToken(RefreshCookie.Value)
}

func setAuthandRefreshCookies(w *http.ResponseWriter, authTokenString string, refreshTokenString string) {
	authCookie := http.Cookie{
		Name:     "AuthToken",
		Value:    authTokenString,
		HttpOnly: true,
	}

	http.SetCookie(*w, &authCookie)

	refreshCookie := http.Cookie{
		Name:     "RefreshToken",
		Value:    refreshTokenString,
		HttpOnly: true,
	}

	http.SetCookie(*w, &refreshCookie)
}

func getCsrfToken(r *http.Request) string {
	csrfFromFrom := r.FormValue("X-CSRF-Token")

	if csrfFromFrom != "" {
		return csrfFromFrom
	} else {
		return r.Header.Get("X-CSRF-Token")
	}
}
