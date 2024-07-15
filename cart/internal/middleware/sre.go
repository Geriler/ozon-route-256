package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestHistogram = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "request_duration_seconds",
		Help:      "Статистика длительности запросов к сервису",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"url", "status"},
)

func ObserveRequestDurationSeconds(duration float64, url, status string) {
	requestHistogram.With(prometheus.Labels{"url": url, "status": status}).Observe(duration)
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
		ObserveRequestDurationSeconds(time.Since(createdAt).Seconds(), fmt.Sprintf("%s %s", r.Method, r.URL.String()), strconv.Itoa(lrw.statusCode))
	}(time.Now())

	h.wrap.ServeHTTP(lrw, r)
}
