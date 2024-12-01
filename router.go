package servant

import (
	"net/http"

	"github.com/gregoryv/htlog"
)

func NewRouter(sys *System) http.HandlerFunc {
	mx := http.NewServeMux()
	mx.Handle("/{$}", home(sys))
	mx.Handle("/login", login(sys))
	mx.Handle("/enter", enter(sys))
	// reuse the same callback endpoint
	mx.Handle("/oauth/redirect", callback(sys))
	mx.Handle("/static/", http.FileServerFS(asset))

	prv := private(mx, sys)
	prv("/inside", inside(sys))
	prv("/settings", settings(sys))

	return logRequests(mx)
}

func home(sys *System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := NewViewModel(sys)
		m.SetSession(sys.ExistingSession(r))
		htdocs.ExecuteTemplate(w, "index.html", m)
	}
}

func login(sys *System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := NewViewModel(sys)
		m.DecorateLogins(r.URL.Query().Get("dest"))
		htdocs.ExecuteTemplate(w, "login.html", m)
	}
}

func enter(sys *System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guardname := r.URL.Query().Get("use")
		// destination after authorized
		dest := r.URL.Query().Get("dest")
		if dest == "" {
			dest = "/inside"
		}
		url, err := sys.GuardURL(guardname, dest)
		if err != nil {
			debug.Printf("login: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func callback(sys *System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		s, err := sys.Authorize(ctx, r)
		if err != nil {
			debug.Printf("callback: %v", err)
			htdocs.ExecuteTemplate(w, "error.html", err)
			return
		}

		// return a page just to set a cookie and then redirect to a
		// location. Cannot set a cookie in a plain redirect response.
		cookie := NewCookie(s.State)
		http.SetCookie(w, cookie)
		m := map[string]string{
			"Location": s.dest, // default page after login
		}

		htdocs.ExecuteTemplate(w, "redirect.html", m)
	}
}

// once authorized, the Contact is inside
func inside(sys *System) privateFunc {
	return func(w http.ResponseWriter, r *http.Request, s *Session) {
		m := NewViewModel(sys)
		m.SetSession(s)
		htdocs.ExecuteTemplate(w, "inside.html", m)
	}
}

func settings(sys *System) privateFunc {
	return func(w http.ResponseWriter, r *http.Request, s *Session) {
		m := NewViewModel(sys)
		m.SetSession(s)
		htdocs.ExecuteTemplate(w, "settings.html", m)
	}
}

func private(mx *http.ServeMux, sys *System) func(string, privateFunc) {
	return func(pattern string, next privateFunc) {
		mx.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			if err := sys.SessionValid(r); err != nil {
				m := map[string]string{
					// page where user selects login
					"Location": "/login?dest=" + r.URL.String(),
				}
				htdocs.ExecuteTemplate(w, "redirect.html", m)
				return
			}
			s := sys.ExistingSession(r)
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

func NewViewModel(sys *System) *ViewModel {
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
		// should match what ever is configured in the system
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

// decorate login links with destination
func (m *ViewModel) DecorateLogins(v string) {
	if v == "" {
		return
	}
	for i, _ := range m.Logins {
		m.Logins[i].Href += "&dest=" + v
	}
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
