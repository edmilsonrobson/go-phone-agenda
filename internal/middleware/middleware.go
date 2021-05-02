package middlewares

import (
	"net/http"
	"time"

	"github.com/edmilsonrobson/go-phone-agenda/internal/logs"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		defer func() {
			logs.InfoLogger.Printf("%v %v in %v\n", r.Method, r.RequestURI, time.Since(t1))
		}()

		next.ServeHTTP(w, r)
	})
}
