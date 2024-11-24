package servant

import (
	"net/http"

	"github.com/gregoryv/servant/htsec"
	"golang.org/x/oauth2"
)

func authLayer(sec *htsec.Secure, next http.Handler) *http.ServeMux {
	mx := http.NewServeMux()
	// explicitly set public patterns so that we don't accidently
	// forget to protect a new endpoint
	mx.Handle("/login", login(sec))
	// reuse the samme callback endpoint
	mx.Handle("/oauth/redirect", callback(sec))
	mx.Handle("/{$}", next)

	// everything else is private
	mx.Handle("/", protect(next))
	return mx
}

func login(sec *htsec.Secure) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		use := r.URL.Query().Get("use")
		Auth, err := sec.AuthService(use)
		if err != nil {
			debug.Printf("login: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		url := Auth.AuthCodeURL(newState(use))
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func protect(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := sessionValid(r); err != nil {
			debug.Printf("protect: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		next.ServeHTTP(w, r)
	}
}

func callback(sec *htsec.Secure) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if err := verify(state); err != nil {
			debug.Printf("callback: %v", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// which Auth service was used
		Auth, err := sec.AuthService(parseUse(state))
		if err != nil {
			debug.Printf("callback: %v", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// get the token
		ctx := oauth2.NoContext
		token, err := Auth.Exchange(ctx, r.FormValue("code"))
		if err != nil {
			debug.Printf("callback oauth exchange: %v", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// get user information from the Auth service
		user, err := Auth.ReadUser(token)
		if err != nil {
			debug.Printf("callback readUser: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newSession(token, user)

		// return a page just to set a cookie and then redirect to a
		// location. Cannot set a cookie in a plain redirect response.
		cookie := newCookie(token)
		http.SetCookie(w, cookie)
		m := map[string]string{
			"Location": "/inside",
		}
		htdocs.ExecuteTemplate(w, "redirect.html", m)
	}
}
