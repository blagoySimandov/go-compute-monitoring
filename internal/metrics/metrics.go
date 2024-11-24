package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Metrics struct {
	ComputeDuration prometheus.Histogram
	ComputeOps      prometheus.Counter
	CPUUsage        prometheus.Gauge
	MemoryUsage     prometheus.Gauge
	GoroutinesCount prometheus.Gauge
	MatrixGenTime   prometheus.Histogram
}

func NewMetrics() *Metrics {
	return &Metrics{
		ComputeDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "matrix_compute_duration_seconds",
			Help:    "Time taken to perform matrix multiplication",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		}),
		ComputeOps: promauto.NewCounter(prometheus.CounterOpts{
			Name: "matrix_compute_operations_total",
			Help: "Total number of matrix multiplications performed",
		}),
		CPUUsage: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "Current CPU usage percentage",
		}),
		MemoryUsage: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Current memory usage in bytes",
		}),
		GoroutinesCount: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "goroutines_count",
			Help: "Current number of goroutines",
		}),
		MatrixGenTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "matrix_generation_duration_seconds",
			Help:    "Time taken to generate matrices",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		}),
	}
}

func (m *Metrics) RecordMetrics() {
	go func() {
		for {
			if cpuPercent, err := cpu.Percent(0, false); err == nil && len(cpuPercent) > 0 {
				m.CPUUsage.Set(cpuPercent[0])
			}

			if memStats, err := mem.VirtualMemory(); err == nil {
				m.MemoryUsage.Set(float64(memStats.Used))
			}

			m.GoroutinesCount.Set(float64(runtime.NumGoroutine()))

			time.Sleep(1 * time.Second)
		}
	}()
}
