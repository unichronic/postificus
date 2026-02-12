package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	// Storage Metrics
	DBQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "db_query_duration_seconds",
		Help:    "Duration of Database queries in seconds",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"operation", "table"})

	RedisOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "redis_operation_duration_seconds",
		Help:    "Duration of Redis operations in seconds",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
	}, []string{"operation"})

	// Business Metrics
	PostPublishTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "post_publish_total",
		Help: "Total number of post publish attempts",
	}, []string{"platform", "status"})

	PostPublishDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "post_publish_duration_seconds",
		Help:    "Duration of post publishing in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"platform"})

	CircuitBreakerState = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "circuit_breaker_state",
		Help: "State of the circuit breaker (0=Closed, 1=Open, 2=HalfOpen)",
	}, []string{"name"})

	// Worker Metrics
	WorkerJobDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "worker_job_duration_seconds",
		Help:    "Duration of worker jobs in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"job_type"})

	WorkerJobsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "worker_jobs_processed_total",
		Help: "Total number of worker jobs processed",
	}, []string{"job_type", "status"})
)
