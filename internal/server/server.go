package server

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())

	http.HandleFunc("/compute", func(w http.ResponseWriter, r *http.Request) {
		compute.PerformComputation(s.metrics.ComputeDuration, s.metrics.ComputeOps, s.metrics.MatrixGenTime)
		w.Write([]byte("Computation completed"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/spawn", func(w http.ResponseWriter, r *http.Request) {
		go s.backgroundComputation(ctx)
		w.Write([]byte("Background computation started"))
	})
	http.HandleFunc("/killall", func(w http.ResponseWriter, r *http.Request) {
		cancel()
		for ctx.Err() == nil {
			// wait for the context to be canceled
		}
		w.Write([]byte("Killed all"))
	})

	go s.backgroundComputation(ctx)

	log.Println("Server starting on :8080")
	log.Printf("Matrix size: %dx%d\n", compute.MatrixSize, compute.MatrixSize)
	return http.ListenAndServe(":8080", nil)
}

func (s *Server) backgroundComputation(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			log.Println("Context canceled, stopping background computation")
			break
		}
		compute.PerformComputation(s.metrics.ComputeDuration, s.metrics.ComputeOps, s.metrics.MatrixGenTime)
		time.Sleep(2 * time.Second)
	}
}
