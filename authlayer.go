package servant

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gregoryv/oauth"
)

func authLayer(next http.Handler) *http.ServeMux {
	mx := http.NewServeMux()
	// explicitly set public patterns so that we don't accidently
	// forget to protect a new endpoint
	mx.Handle("/login", github.Login())
	// todo github is just one of the available auth sources
	mx.Handle("GET "+github.RedirectPath(), github.Authorize(enter))
	mx.Handle("/{$}", next)

	// everything else is private
	mx.Handle("/", protect(next))
	return mx
}

func protect(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		_, sessionFound := sessions[token.Value]
		if err != nil || !sessionFound {
			debug.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// enter is used after a user authenticates via github. It sets a
// token cookie.
func enter(token string, w http.ResponseWriter, r *http.Request) {
	var user struct {
		Email string
		Name  string
	}

	resp, err := http.DefaultClient.Do(github.User(token))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewDecoder(resp.Body).Decode(&user)

	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	}
	// cache the session
	session := Session{
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	}
	sessions[session.Token] = session
	debug.Println(session.String())
	// return a page just to set a cookie and then redirect to a
	// location. Cannot set a cookie in a plain redirect response.
	http.SetCookie(w, &cookie)
	m := map[string]string{
		"Location": "/inside",
	}
	page.ExecuteTemplate(w, "redirect.html", m)
}

// use x/oauth2 with github.Endpoint
var github = oauth.Github{
	ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
	ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
	RedirectURI:  os.Getenv("OAUTH_GITHUB_REDIRECT_URI"),
}
