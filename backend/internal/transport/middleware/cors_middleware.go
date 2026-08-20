package middleware

import "net/http"

type CORSMiddleware struct {
	allowedHost string
}

func NewCORSMiddleware(allowedHost string) CORSMiddleware {
	return CORSMiddleware{
		allowedHost: allowedHost,
	}
}

func (m CORSMiddleware) ActivateCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", m.allowedHost)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
