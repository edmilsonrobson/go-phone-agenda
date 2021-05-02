package middlewares

import (
	"net/http"
	"time"

	"github.com/edmilsonrobson/go-phone-agenda/internal/logs"
)

type responseObserver struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		observerWriter := &responseObserver{ResponseWriter: w}
		defer func() {
			logs.InfoLogger.Printf("(HTTP %v) %v %v in %v\n", observerWriter.status, r.Method, r.RequestURI, time.Since(t1))
		}()

		next.ServeHTTP(observerWriter, r)
	})
}

func (w *responseObserver) Status() int {
	return w.status
}

func (w *responseObserver) Write(p []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func (w *responseObserver) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	// Check after in case there's error handling in the wrapped ResponseWriter.
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
}
