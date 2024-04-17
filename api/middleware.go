package api

import (
	"log"
	"net/http"
	"ticketing-api/auth"
	"ticketing-api/types"
	"time"
)

func IsAdmin(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := auth.IsRole(r, types.RoleAdmin)
		if err != nil {
			encodeResponse(w, http.StatusUnauthorized, &APIResponse{Status: http.StatusUnauthorized, Message: err.Error()})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IsEditor(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := auth.IsRole(r, types.RoleEditor)
		if err != nil {
			encodeResponse(w, http.StatusUnauthorized, &APIResponse{Status: http.StatusUnauthorized, Message: err.Error()})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IsAuthenticated(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := auth.IsAuthenticated(r)
		if err != nil {
			encodeResponse(w, http.StatusUnauthorized, &APIResponse{Status: http.StatusUnauthorized, Message: err.Error()})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

type middleware func(http.Handler) http.HandlerFunc

func CreateStack(mw ...middleware) middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(mw) - 1; i >= 0; i-- {
			next = mw[i](next)
		}

		return next.ServeHTTP
	}
}
