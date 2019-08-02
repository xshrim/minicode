package sdk

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		rawuri, _ := url.QueryUnescape(r.RequestURI)

		log.Printf(
			"%s [%s] %s %s",
			r.Method,
			rawuri,
			name,
			time.Since(start),
		)
	})
}
