// midleware chain , for TODO some common operation in the future before request processing
package webapp

import (
	"net/http"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//change context
		//TODO smth

		newR := r //remove
		next.ServeHTTP(w, newR)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//TODO smth
		next.ServeHTTP(w, r)
	})
}

func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		next.ServeHTTP(w, r)
	})
}
