package gol

import (
	"net/http"
	"time"
)

// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields

type logRespWriter struct {
	http.ResponseWriter
	statusCode int
}

// wrap response writer with status code support
func NewLogRespWriter(w http.ResponseWriter) *logRespWriter {
	return &logRespWriter{w, http.StatusOK}
}

func (lrw *logRespWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// wrap http handler function with log support
func HttpHandlerFunc(fn http.HandlerFunc, headers ...string) http.HandlerFunc {
	return std.HttpHandlerFunc(fn, headers...)
}

// wrap http handler with log support
func HttpHandler(handler http.Handler, headers ...string) http.Handler {
	return std.HttpHandler(handler, headers...)
}

// wrap http handler function with log support
func (l *Logger) HttpHandlerFunc(fn http.HandlerFunc, headers ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		wl := NewLogRespWriter(w)
		fn(wl, r)

		if l.HasFlag(Ljson) {
			data := map[string]interface{}{
				"status-code":    wl.statusCode,
				"latency-time":   time.Since(startTime),
				"client-ip":      r.RemoteAddr,
				"request-method": r.Method,
				"request-uri":    r.RequestURI,
			}

			for _, hd := range headers {
				data[toLower(hd)] = r.Header.Get(hd)
			}

			l.Info(string(map2json(nil, data)))

		} else if l.HasFlag(Lcolor) && !l.HasFlag(Lfcolor) {
			format := "| %s | %10v | %15s | %s | %s |"
			values := []interface{}{
				colorStatusCode(wl.statusCode),
				time.Since(startTime),
				r.RemoteAddr,
				colorRequestMethod(r.Method),
				r.RequestURI,
			}

			for _, hd := range headers {
				format += " %v |"
				values = append(values, r.Header.Get(hd))
			}

			l.Infof(format, values...)

		} else {
			format := "| %3d | %10v | %15s | %5s | %s |"
			values := []interface{}{
				wl.statusCode,
				time.Since(startTime),
				r.RemoteAddr,
				r.Method,
				r.RequestURI,
			}

			for _, hd := range headers {
				format += " %v |"
				values = append(values, r.Header.Get(hd))
			}

			l.Infof(format, values...)
		}
	}
}

// wrap http handler with log support
func (l *Logger) HttpHandler(handler http.Handler, headers ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		wl := NewLogRespWriter(w)
		handler.ServeHTTP(wl, r)

		if l.HasFlag(Ljson) {
			data := map[string]interface{}{
				"status-code":    wl.statusCode,
				"latency-time":   time.Since(startTime),
				"client-ip":      r.RemoteAddr,
				"request-method": r.Method,
				"request-uri":    r.RequestURI,
			}

			for _, hd := range headers {
				data[toLower(hd)] = r.Header.Get(hd)
			}

			l.Info(string(map2json(nil, data)))

		} else if l.HasFlag(Lcolor) && !l.HasFlag(Lfcolor) {
			format := "| %s | %10v | %15s | %s | %s |"
			values := []interface{}{
				colorStatusCode(wl.statusCode),
				time.Since(startTime),
				r.RemoteAddr,
				colorRequestMethod(r.Method),
				r.RequestURI,
			}

			for _, hd := range headers {
				format += " %v |"
				values = append(values, r.Header.Get(hd))
			}

			l.Infof(format, values...)

		} else {
			format := "| %3d | %10v | %15s | %5s | %s |"
			values := []interface{}{
				wl.statusCode,
				time.Since(startTime),
				r.RemoteAddr,
				r.Method,
				r.RequestURI,
			}

			for _, hd := range headers {
				format += " %v |"
				values = append(values, r.Header.Get(hd))
			}

			l.Infof(format, values...)
		}
	})
}
