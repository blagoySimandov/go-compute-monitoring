package server

import (
	"log"
	"matrix-compute/internal/compute"
	"matrix-compute/internal/metrics"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	metrics *metrics.Metrics
}

func NewServer() *Server {
	return &Server{
		metrics: metrics.NewMetrics(),
	}
}

func (s *Server) Start() error {
	s.metrics.RecordMetrics()

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/compute", func(w http.ResponseWriter, r *http.Request) {
		compute.PerformComputation(s.metrics.ComputeDuration, s.metrics.ComputeOps, s.metrics.MatrixGenTime)
		w.Write([]byte("Computation completed"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	go s.backgroundComputation()

	log.Println("Server starting on :8080")
	log.Printf("Matrix size: %dx%d\n", compute.MatrixSize, compute.MatrixSize)
	return http.ListenAndServe(":8080", nil)
}

func (s *Server) backgroundComputation() {
	for {
		compute.PerformComputation(s.metrics.ComputeDuration, s.metrics.ComputeOps, s.metrics.MatrixGenTime)
		time.Sleep(2 * time.Second)
	}
}
