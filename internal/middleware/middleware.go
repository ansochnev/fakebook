package middleware

import (
	"net/http"
	"strings"
)

func RemoveTrailingSlashFromPath(next http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			req.URL.Path = strings.TrimRight(req.URL.Path, "/")
		}
		next.ServeHTTP(rw, req)
	}
}
