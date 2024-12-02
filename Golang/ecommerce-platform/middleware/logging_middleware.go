// // middleware/logging_middleware.go
// package middleware

// import (
//     "net/http"
//     "time"
//     "log"
// )

// func LoggingMiddleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         start := time.Now()
//         log.Printf("Request: %s %s", r.Method, r.URL.Path)
//         next.ServeHTTP(w, r)
//         duration := time.Since(start)
//         log.Printf("Request: %s %s %s", r.Method, r.URL.Path, duration)
//     })
// }


// type statusRecorder struct {
// 	http.ResponseWriter
// 	statusCode int
// }

// func (r *statusRecorder) WriteHeader(code int) {
// 	r.statusCode = code
// 	r.ResponseWriter.WriteHeader(code)
// }

package middleware

import (
	"net/http"
	"time"

	"github.com/Tasherokk/GolangProject.git/utils" // Adjust the import path to your project
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start the timer
		start := time.Now()

		// Wrap the ResponseWriter to capture the status code
		statusRecorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// Log the request
		next.ServeHTTP(statusRecorder, r)

		// Calculate the duration
		duration := time.Since(start).Seconds()

		// Update Prometheus metrics
		utils.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusRecorder.statusCode)).Inc()
		utils.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
