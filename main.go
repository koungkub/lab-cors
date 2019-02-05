package main

import (
	"net/http"
)

func main() {
	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("login")
		if err != nil || c.Value != "1" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`supersecret response`))
	})

	m.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Path:     "/",
			Name:     "login",
			Value:    "1",
			HttpOnly: true,
			MaxAge:   3600,
		})
		w.Write([]byte("login success"))
	})

	m.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Path:     "/",
			Name:     "login",
			HttpOnly: true,
			MaxAge:   -1,
		})
		w.Write([]byte("logout success"))
	})

	http.ListenAndServe(":8080", middleware(m))
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// iframe
		w.Header().Set("X-Frame-Options", "deny")

		// check origin
		if origin := r.Header.Get("Origin"); len(origin) > 0 {
			if origin != "http://localhost:3000" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		// cors
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "5")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// ajax
		if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}
