package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestHistogramDatabase = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "request_database_duration_seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	requestHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "request_duration_seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"url", "status"},
	)
)

func SaveDatabaseMetrics(duration float64, method, status string) {
	requestHistogramDatabase.With(prometheus.Labels{"method": method, "status": status}).Observe(duration)
}

func SaveLomsMetrics(duration float64, url, status string) {
	requestHistogram.With(prometheus.Labels{"url": url, "status": status}).Observe(duration)
}
