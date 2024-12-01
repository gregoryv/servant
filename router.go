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
	mx.Handle("/enter", enter(sec))
	// reuse the same callback endpoint
	mx.Handle("/oauth/redirect", callback(sec))
	mx.Handle("/static/", http.FileServerFS(asset))

	prv := private(mx)
	prv("/inside", inside)
	prv("/settings", settings)

	return logRequests(mx)
}

// todo use a view model

func frontpage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htdocs.ExecuteTemplate(w, "index.html", nil)
	}
}

func login(sec *htsec.Detail) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]string{
			"PathLoginGithub": "/enter?use=github",
			"PathLoginGoogle": "/enter?use=google",
		}
		if v := r.URL.Query().Get("dest"); v != "" {
			for k, _ := range m {
				m[k] += "&dest=" + v
			}
		}
		htdocs.ExecuteTemplate(w, "login.html", m)
	}
}

func enter(sec *htsec.Detail) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guardname := r.URL.Query().Get("use")
		// destination after authorized
		dest := r.URL.Query().Get("dest")
		if dest == "" {
			dest = "/inside"
		}
		url, err := sec.GuardURL(guardname, dest)
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
			"Location": slip.Dest(), // default page after login
		}

		htdocs.ExecuteTemplate(w, "redirect.html", m)
	}
}

// once authorized, the Contact is inside
func inside(w http.ResponseWriter, r *http.Request, s *Session) {
	htdocs.ExecuteTemplate(w, "inside.html", s)
}

func settings(w http.ResponseWriter, r *http.Request, s *Session) {
	htdocs.ExecuteTemplate(w, "settings.html", existingSession(r))
}

func private(mx *http.ServeMux) func(string, privateFunc) {
	return func(pattern string, next privateFunc) {
		mx.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			if err := sessionValid(r); err != nil {
				debug.Printf("protect: %v", err)
				m := map[string]string{
					// page where user selects login
					"Location": "/login?dest=" + r.URL.String(),
				}
				htdocs.ExecuteTemplate(w, "redirect.html", m)
				return
			}
			s := existingSession(r)
			next(w, r, &s)
		})
	}
}

type privateFunc func(http.ResponseWriter, *http.Request, *Session)

func logRequests(next http.Handler) http.HandlerFunc {
	log := htlog.Middleware{
		Println: debug.Println,
		Clean:   htlog.QueryHide("access_token"),
	}
	return log.Use(next)
}
