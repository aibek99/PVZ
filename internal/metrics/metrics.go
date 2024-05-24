package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// TotalRequests is
	TotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests",
	})

	// ResponseTimes is
	ResponseTimes = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "grpc_response_times_seconds",
		Help:    "Response times for gRPC requests",
		Buckets: prometheus.LinearBuckets(0.01, 0.05, 20), // Adjust bucket sizes as appropriate
	})

	// ResponseStatus is
	ResponseStatus = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_response_status",
		Help: "Status of gRPC responses",
	}, []string{"status"})

	// ErrorRates is
	ErrorRates = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_error_rates_total",
		Help: "Total number of error responses in gRPC requests",
	})
)

// Listen is
func Listen(metricsPort string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         metricsPort,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	return server.ListenAndServe()
}
