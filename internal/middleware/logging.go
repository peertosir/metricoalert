package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		loggingRespWriter := middleware.NewWrapResponseWriter(rw, req.ProtoMajor)
		handler.ServeHTTP(loggingRespWriter, req)
		reqDuration := time.Since(startTime)
		responseStatusCode := loggingRespWriter.Status()
		// dont know why in case of POST request status code is 0, this is workaround
		if responseStatusCode == 0 {
			responseStatusCode = http.StatusOK
		}
		log.Info().Msgf(
			"Request: [%s]%s [%d] time=%dns response size = %d",
			req.Method, req.URL, responseStatusCode,
			reqDuration.Nanoseconds(), loggingRespWriter.BytesWritten(),
		)
	},
	)
}
