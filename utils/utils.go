package utils

import (
	"log"
	"net/http"
	"runtime"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		log.Printf(
			`{"proto": "%s", "route" :"%s%s", "method": "%s", "time": "%s"}`,
			r.Proto, r.Host, r.RequestURI, r.Method, time.Since(t),
		)
		next.ServeHTTP(w, r)
	})
}

func ErrorLineLogger(err error) {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("[error] %s:%d %s", file, line, err)
}
