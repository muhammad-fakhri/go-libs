package corsanywhere

import (
	"net/http"
)

// CorsAnywhere is the main middleware function
// it adds Origin, Access-Control-Allow-Origin, Access-Control-Allow-Credentials to response headers
// and additionally handle prefetch cors request
func CorsAnywhere(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := "*"
		if r.Header.Get("Origin") != "" {
			origin = r.Header.Get("Origin")
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)

		if rm := r.Header.Get("Access-Control-Request-Method"); rm != "" {
			w.Header().Set("Access-Control-Allow-Methods", rm)
		}

		if rh := r.Header.Get("Access-Control-Request-Headers"); rh != "" {
			w.Header().Set("Access-Control-Allow-Headers", rh)
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// not prefetch request thus we forward it to the next handler
		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)

		w.Write(nil)
	})
}
