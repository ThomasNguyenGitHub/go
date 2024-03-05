package middleware

import (
	"github.com/ThomasNguyenGitHub/go/log"
	"net/http"
	"runtime/debug"
)

func (m *middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("err:", err, "stack-trace", string(debug.Stack()))

			}
		}()
		next.ServeHTTP(w, r)
	})
}
