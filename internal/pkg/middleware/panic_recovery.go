package middleware

import (
	"log"
	"net/http"
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/pkg/response"
	"runtime/debug"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v\n%s", err, debug.Stack())
				response.SendErrorResponse(w, api.INTERNAL, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
