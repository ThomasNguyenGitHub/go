package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/ThomasNguyenGitHub/go/cache"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Middleware interface {
	Recover(next http.Handler) http.Handler
}

type middleware struct {
	redisCache cache.Cacher
}

func NewMiddleware(redisCache cache.Cacher) Middleware {
	return &middleware{
		redisCache: redisCache,
	}
}

func marshall(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func combineMethodAndPath(r *http.Request) string {
	var (
		method = r.Method
		ep, _  = mux.CurrentRoute(r).GetPathTemplate()
		idx    = strings.LastIndex(ep, "/")
	)
	if idx >= 0 {
		ep = ep[idx+1:]
	}
	return fmt.Sprintf("%s_%s", strings.ToLower(method), strings.ReplaceAll(ep, "-", "_"))
}
