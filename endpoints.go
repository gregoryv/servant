package servant

import "net/http"

func newRouter() http.Handler {
	mx := http.NewServeMux()
	// any Auth related endpoints are defined in the AuthLayer
	mx.Handle("/{$}", frontpage())
	mx.Handle("/inside", inside())
	mx.Handle("/settings", settings())
	return mx
}

func frontpage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]any{
			"PathLoginGithub": "/login?use=github",
		}
		htdocs.ExecuteTemplate(w, "index.html", m)
	}
}

// once authenticated, the user is inside
func inside() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htdocs.ExecuteTemplate(w, "inside.html", existingSession(r))
	}
}

func settings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htdocs.ExecuteTemplate(w, "settings.html", existingSession(r))
	}
}
