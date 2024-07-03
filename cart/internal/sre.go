package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestHistogram = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "request_duration_seconds",
		Help:      "Duration of requests in seconds.",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"url", "status"},
)

func SaveMetrics(duration float64, url, status string) {
	requestHistogram.With(prometheus.Labels{"url": url, "status": status}).Observe(duration)
}
