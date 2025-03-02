package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request count
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Kafka event processing latency
	KafkaProcessingLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_processing_latency_seconds",
			Help:    "Latency of Kafka event processing",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"},
	)

	// Click event count
	ClickEventsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "click_events_total",
			Help: "Total number of click events processed",
		},
	)
)
