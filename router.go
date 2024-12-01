package servant

import (
	"net/http"

	"github.com/gregoryv/htlog"
	"github.com/gregoryv/htsec"
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

	return logRequests(mx)
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

func login(sec *htsec.Detail) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guardname := r.URL.Query().Get("use")
		url, err := sec.GuardURL(guardname)
		if err != nil {
			debug.Printf("login: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func callback(sec *htsec.Detail) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slip, err := sec.Authorize(ctx, r)
		if err != nil {
			debug.Printf("callback: %v", err)
			htdocs.ExecuteTemplate(w, "error.html", err)
			return
		}

		newSession(slip)

		// return a page just to set a cookie and then redirect to a
		// location. Cannot set a cookie in a plain redirect response.
		cookie := newCookie(slip.State)
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

// once authorized, the Contact is inside
func inside(w http.ResponseWriter, r *http.Request, s *Session) {
	htdocs.ExecuteTemplate(w, "inside.html", s)
}

func settings(w http.ResponseWriter, r *http.Request, s *Session) {
	htdocs.ExecuteTemplate(w, "settings.html", existingSession(r))
}

func logRequests(next http.Handler) http.HandlerFunc {
	log := htlog.Middleware{
		Println: debug.Println,
		Clean:   htlog.QueryHide("access_token"),
	}
	return log.Use(next)
}
