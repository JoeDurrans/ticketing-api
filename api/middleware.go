package api

import (
	"log"
	"net/http"
	"ticketing-api/auth"
	"time"
)

func IsAdmin(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok := auth.IsAdmin(r); ok {
			next.ServeHTTP(w, r)
			return
		}

		EncodeResponse(w, http.StatusUnauthorized, &APIResponse{Status: http.StatusUnauthorized, Message: "permission denied"})
	})
}

func IsEditor(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok := auth.IsEditor(r); ok {
			next.ServeHTTP(w, r)
			return
		}

		EncodeResponse(w, http.StatusUnauthorized, &APIResponse{Status: http.StatusUnauthorized, Message: "permission denied"})
	})
}

func IsAuthenticated(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok := auth.IsAuthenticated(r); ok {
			next.ServeHTTP(w, r)
			return
		}

		EncodeResponse(w, http.StatusUnauthorized, &APIResponse{Status: http.StatusUnauthorized, Message: "permission denied"})
	})
}

func Logging(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

type Middleware func(http.Handler) http.HandlerFunc

func CreateStack(mw ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(mw) - 1; i >= 0; i-- {
			next = mw[i](next)
		}

		return next.ServeHTTP
	}
}
