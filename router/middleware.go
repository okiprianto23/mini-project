package router

import (
	"net/http"
)

func MiddlewareCustom(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//now := time.Now()

			CORSOriginHandler(&w)

			if r.Method != http.MethodOptions {
				next.ServeHTTP(w, r)
			}

			//logger.With(zap.Int64("processing_time", time.Since(now).Microseconds()))

			//statusCode := 200
			//
			//if r.URL.Path != "/metrics" && r.Method != http.MethodOptions {
			//	if m.metrics != nil {
			//		usedPath, _ := mux.CurrentRoute(r).GetPathTemplate()
			//		if usedPath == "" {
			//			usedPath = r.URL.Path
			//		}
			//
			//		m.metrics.GetDefaultMetric().APIHist.WithLabelValues(
			//			usedPath,
			//			r.Method,
			//			strconv.Itoa(statusCode),
			//		).Observe(float64(time.Since(now).Seconds()))
			//	}
			//}
		},
	)
}

func CORSOriginHandler(responseWriter *http.ResponseWriter) {
	(*responseWriter).Header().Set("Access-Control-Allow-Origin", "*")
	(*responseWriter).Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, authorization, x-nextoken, content-length")
	(*responseWriter).Header().Set("Access-Control-Allow-Credentials", "true")
	(*responseWriter).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	(*responseWriter).Header().Set("Access-Control-Max-Age", "1209600")
}
