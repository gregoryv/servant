package servant

import "net/http"

func endpoints() http.Handler {
	mx := http.NewServeMux()
	// any auth related endpoints are defined in the AuthLayer
	mx.Handle("/{$}", frontpage())
	mx.Handle("/inside", inside())
	mx.Handle("/settings", settings())
	return mx
}

func frontpage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]any{
			"PathLoginGithub": "/login",
		}
		page.ExecuteTemplate(w, "index.html", m)
	}
}

// once authenticated, the user is inside
func inside() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page.ExecuteTemplate(w, "inside.html", existingSession(r))
	}
}

func settings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page.ExecuteTemplate(w, "settings.html", existingSession(r))
	}
}
