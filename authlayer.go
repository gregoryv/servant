package servant

import (
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

func authLayer(next http.Handler) *http.ServeMux {
	mx := http.NewServeMux()
	// explicitly set public patterns so that we don't accidently
	// forget to protect a new endpoint
	state := "RANDOM_STATE"
	mx.Handle("/login", login(state))
	// todo github is just one of the available auth sources
	mx.Handle("/oauth/redirect", callback(state))
	mx.Handle("/{$}", next)

	// everything else is private
	mx.Handle("/", protect(next))
	return mx
}

func login(state string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := githubOauth.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func protect(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			debug.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, sessionFound := sessions[token.Value]
		if !sessionFound {
			debug.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func callback(state string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != state {
			debug.Print("missing state")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		ctx := oauth2.NoContext
		token, err := githubOauth.Exchange(ctx, r.FormValue("code"))
		if err != nil {
			debug.Print("oauth exchange:", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		{
			r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
			r.Header.Set("Accept", "application/vnd.github.v3+json")
			r.Header.Set("Authorization", "token "+token.AccessToken)
			resp, err := http.DefaultClient.Do(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var u user
			json.NewDecoder(resp.Body).Decode(&u)
			newSession(token, &u)
		}

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

var githubOauth = &oauth2.Config{
	RedirectURL:  os.Getenv("OAUTH_GITHUB_REDIRECT_URI"),
	ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
	ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
	Endpoint:     endpoints.GitHub,
}
