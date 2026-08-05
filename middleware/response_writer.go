package middleware

import "net/http"

/*
responseWriter wraps http.ResponseWriter so response metadata
can be captured by middleware.
*/
type responseWriter struct {
	http.ResponseWriter

	statusCode int
	upstreamID string
}

/*
WriteHeader captures the status code before forwarding it
to the underlying ResponseWriter.
*/
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

/*
Write ensures a default HTTP 200 status code is recorded
when the downstream handler writes a response body without
explicitly calling WriteHeader.
*/
func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(data)
}

//SetUpstream records the upstream that handled the request.
func (rw *responseWriter) SetUpstream(id string) {
	rw.upstreamID = id
}
