package middleware

import (
	"net/http"
)

func CORSMethodMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
			w.Header().Set("Access-Control-Allow-Headers", "client_ip, user_agent, x_forwarded_for, date, "+
				"x-amz-apigw-id, x-amzn-errortype, x-amzn-requestid, Accept, Accept-Language, Content-Type, Authorization, x-ijt, *")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
