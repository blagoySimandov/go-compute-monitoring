package server

import (
	"context"
	"log"
	"matrix-compute/internal/compute"
	"matrix-compute/internal/metrics"
	"net/http"
	"runtime"
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
		compute.PerformComputation(s.metrics.ComputeDuration, s.metrics.ComputeOps, s.metrics.MatrixGenTime, s.metrics.MatrixSize, s.metrics.MatrixCount)
		w.Write([]byte("Computation completed"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/spawn", func(w http.ResponseWriter, r *http.Request) {
		// Spawn multiple background computations based on CPU count
		for i := 0; i < runtime.NumCPU(); i++ {
			go s.backgroundComputation(ctx)
		}
		w.Write([]byte("Background computations started"))
	})

	http.HandleFunc("/killall", func(w http.ResponseWriter, r *http.Request) {
		cancel()
		for ctx.Err() == nil {
			// wait for the context to be canceled
		}
		w.Write([]byte("Killed all"))
	})
	for i := 0; i < runtime.NumCPU(); i++ {
		go s.backgroundComputation(ctx)
	}

	stop := compute.GraduallyIncreaseMatrixSize()
	http.HandleFunc("/stop-load-increase", func(w http.ResponseWriter, r *http.Request) {
		stop()
		w.Write([]byte("stopped the load increase"))
	})

	log.Println("Server starting on :8080")
	log.Printf("Initial matrix size: 100x100, will grow up to 1000x1000 (increasing by %d each time)\n", 10)
	log.Printf("Running %d parallel computation routines\n", runtime.NumCPU())
	return http.ListenAndServe(":8080", nil)
}

func (s *Server) backgroundComputation(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			log.Println("Context canceled, stopping background computation")
			break
		}
		compute.PerformComputation(s.metrics.ComputeDuration, s.metrics.ComputeOps, s.metrics.MatrixGenTime, s.metrics.MatrixSize, s.metrics.MatrixCount)
		time.Sleep(1 * time.Second) // Reduced sleep time for more frequent computations
	}
}
