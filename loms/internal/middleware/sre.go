package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	requestDatabaseDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "request_database_duration_seconds",
			Help:      "Статистика длительности запросов к базе данных",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	requestGRPCDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "request_grpc_duration_seconds",
			Help:      "Статистика длительности запросов к сервису по gRPC",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"url", "status"},
	)

	requestHTTPDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "request_http_duration_seconds",
			Help:      "Статистика длительности запросов к сервису по HTTP",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"url", "status"},
	)
)

func ObserveRequestDatabaseDurationSeconds(duration float64, method, status string) {
	requestDatabaseDurationSeconds.With(prometheus.Labels{"method": method, "status": status}).Observe(duration)
}

func ObserveRequestHTTPDurationSeconds(duration float64, url, status string) {
	requestHTTPDurationSeconds.With(prometheus.Labels{"url": url, "status": status}).Observe(duration)
}

func ObserveRequestGRPCDurationSeconds(duration float64, url, status string) {
	requestGRPCDurationSeconds.With(prometheus.Labels{"url": url, "status": status}).Observe(duration)
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

type SreWrapperHandler struct {
	wrap http.Handler
}

func NewSreWrapperHandler(wrap http.Handler) *SreWrapperHandler {
	return &SreWrapperHandler{wrap: wrap}
}

func (h *SreWrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := NewLoggingResponseWriter(w)

	defer func(createdAt time.Time) {
		ObserveRequestHTTPDurationSeconds(time.Since(createdAt).Seconds(), fmt.Sprintf("%s %s", r.Method, r.Pattern), strconv.Itoa(lrw.statusCode))
	}(time.Now())

	h.wrap.ServeHTTP(lrw, r)
}

func GRPCSreWrapper(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	statusCode := "ok"
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			statusCode = s.String()
		}
	}

	ObserveRequestGRPCDurationSeconds(duration.Seconds(), info.FullMethod, statusCode)

	return resp, err
}
