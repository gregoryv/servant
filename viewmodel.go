package servant

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
