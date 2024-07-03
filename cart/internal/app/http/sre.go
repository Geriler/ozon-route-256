package http

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cart",
			Name:      "request_count",
			Help:      "Total number of requests processed.",
		},
		[]string{"url"},
	)

	requestHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "cart",
			Name:      "request_duration_seconds",
			Help:      "Duration of requests in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"url", "status_code"},
	)
)
