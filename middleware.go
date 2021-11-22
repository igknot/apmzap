package logger

import (
	"log"
	"net/http"
	"time"
	"go.elastic.co/apm/module/apmzap"
	"go.uber.org/zap"
)
//StatusRecorder wraps a ResponseWriter and records the status code
type StatusRecorder struct {
	http.ResponseWriter
	Status int
	Length int
	
}
//WriteHeader Intercepts the status
func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
//Write Intercepts the status
func (r *StatusRecorder) Write(b []byte) (n int, err error) {
    n, err = r.ResponseWriter.Write(b)
    r.Length += n

    return
}
//ZapLogger wraps a zap logger 
type zapLogger struct {
	logZ *zap.Logger
	name string
}

// NewMiddleware returns a new Zap Middleware handler.
func NewZapLogger(name string, logger *zap.Logger) func(next http.Handler) http.Handler {
	if logger == nil {
		log.Fatal("No logger provided ")
	}
	return zapLogger{
		logZ: logger,
		name: name,
	}.Middleware
}
//middleware logs tracecontext and request detail to ease correlation with APM trace details 
func (zl zapLogger) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
			Length: 0,
		}
		start := time.Now()
		h.ServeHTTP(recorder, r)
		duration := time.Since(start).Nanoseconds()
		traceContextFields := apmzap.TraceContext(r.Context())
		if traceContextFields == nil {
			traceContextFields = []zap.Field{ 
				zap.String("trace.context", "not-traced"),
			}

		}

		fields := append(
			traceContextFields, 
			zap.String("http.request.method", r.Method),
			zap.String("url.path", r.RequestURI),
			zap.String("logger", zl.name),
			zap.Int("http.Status", recorder.Status),
			zap.Int64("event.duration_nanoseconds", duration),
			zap.Int("http.response.length", recorder.Length),
			zap.String("client.address", r.RemoteAddr),
		)

		if (recorder.Status >= 500 || recorder.Status == 429) {
			zl.logZ.Error("",fields...)
			return
		}
		if recorder.Status >=400   {
			zl.logZ.Warn("",fields...)
			return
		}
		zl.logZ.Info("",fields...)

		
	})
}
