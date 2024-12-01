package servant

import (
	"net/http"

	"github.com/gregoryv/htlog"
	"github.com/gregoryv/htsec"
)

func NewRouter(sys *System) http.HandlerFunc {
	sec := sys.Security()
	mx := http.NewServeMux()
	mx.Handle("/{$}", home())
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

func home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := NewViewModel()
		m.SetSession(existingSession(r))
		htdocs.ExecuteTemplate(w, "index.html", m)
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
	m := NewViewModel()
	m.SetSession(s)
	htdocs.ExecuteTemplate(w, "inside.html", m)
}

func settings(w http.ResponseWriter, r *http.Request, s *Session) {
	m := NewViewModel()
	m.SetSession(s)
	htdocs.ExecuteTemplate(w, "settings.html", m)
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
			next(w, r, s)
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

// ----------------------------------------

func NewViewModel() *ViewModel {
	return &ViewModel{
		Nav: &Nav{
			Home: Link{
				Href: "/",
				Text: "Home",
			},
			Inside: Link{
				Private: true,
				Href:    "/inside",
				Text:    "Inside",
			},
			Settings: Link{
				Private: true,
				Href:    "/settings",
				Text:    "Settings",
			},
			Login: &Link{
				Href: "/login",
				Text: "Login",
			},
		},
		Logins: []GuardLink{
			{
				Img:  "/static/github.svg",
				Href: "/enter?use=github",
				Text: "Github",
			},
			{
				Img:  "/static/google.svg",
				Href: "/enter?use=google",
				Text: "Google",
			},
		},
	}
}

type ViewModel struct {
	*Nav
	Logins  []GuardLink
	Session *Session
}

func (m *ViewModel) SetSession(s *Session) {
	m.Session = s
	m.Nav.SetSession(s)
}

type Nav struct {
	Home     Link
	Inside   Link
	Settings Link
	Login    *Link
}

func (n *Nav) SetSession(s *Session) {
	if s == nil {
		return
	}
	// should it be hidden here? or should the template decide based
	// on session
	n.Login = nil
	n.Inside.Private = false
	n.Settings.Private = false
}

type Link struct {
	Private bool
	Href    string
	Text    string
}

type GuardLink struct {
	Img  string
	Href string
	Text string
}
