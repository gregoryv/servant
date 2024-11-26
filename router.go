package servant

import (
	"net/http"

	"github.com/gregoryv/servant/htsec"
	"golang.org/x/oauth2"
)

func NewRouter(sys *System) http.HandlerFunc {
	sec := sys.Security()
	mx := http.NewServeMux()
	mx.Handle("/{$}", frontpage())
	mx.Handle("/login", login(sec))
	// reuse the same callback endpoint
	mx.Handle("/oauth/redirect", callback(sec))

	// everything else is private
	mx.Handle("/", private())
	return logware(mx)
}

func frontpage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]any{
			"PathLoginGithub": "/login?use=github",
			"PathLoginGoogle": "/login?use=google",
		}
		htdocs.ExecuteTemplate(w, "index.html", m)
	}
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

func callback(sec *htsec.Secure) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if err := verify(state); err != nil {
			debug.Printf("callback: %v", err)
			htdocs.ExecuteTemplate(w, "error.html", err)
			return
		}
		// which Auth service was used
		Auth, err := sec.AuthService(parseUse(state))
		if err != nil {
			debug.Printf("callback: %v", err)
			htdocs.ExecuteTemplate(w, "error.html", err)
			return
		}
		// get the token
		ctx := oauth2.NoContext
		token, err := Auth.Exchange(ctx, r.FormValue("code"))
		if err != nil {
			debug.Printf("callback oauth exchange: %v", err)
			htdocs.ExecuteTemplate(w, "error.html", err)
			return
		}
		// get user information from the Auth service
		user, err := Auth.ReadUser(token)
		if err != nil {
			debug.Printf("callback readUser: %v", err)
			htdocs.ExecuteTemplate(w, "error.html", err)
			return
		}
		newSession(state, token, user)

		// return a page just to set a cookie and then redirect to a
		// location. Cannot set a cookie in a plain redirect response.
		cookie := newCookie(state)
		http.SetCookie(w, cookie)
		m := map[string]string{
			"Location": "/inside",
		}
		htdocs.ExecuteTemplate(w, "redirect.html", m)
	}
}

func private() http.Handler {
	mx := http.NewServeMux()
	handle := withSession(mx)

	handle("/inside", inside)
	handle("/settings", settings)
	return protect(mx)
}

func withSession(mx *http.ServeMux) func(string, privateFunc) {
	return func(pattern string, next privateFunc) {
		mx.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			s := existingSession(r)
			next(w, r, &s)
		})
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

type privateFunc func(http.ResponseWriter, *http.Request, *Session)

// once authenticated, the user is inside
func inside(w http.ResponseWriter, r *http.Request, s *Session) {
	htdocs.ExecuteTemplate(w, "inside.html", s)
}

func settings(w http.ResponseWriter, r *http.Request, s *Session) {
	htdocs.ExecuteTemplate(w, "settings.html", existingSession(r))
}
