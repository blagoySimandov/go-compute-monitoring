package compute

import (
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	baseMatrixSize int32 = 10
	growthRate     int32 = 20
)

func getCurrentMatrixSize() int {
	return int(atomic.LoadInt32(&baseMatrixSize))
}

func GenerateMatrix(size int) [][]float64 {
	matrix := make([][]float64, size)
	for i := range matrix {
		matrix[i] = make([]float64, size)
		for j := range matrix[i] {
			matrix[i][j] = rand.Float64()
		}
	}
	return matrix
}

func MultiplyMatrices(a, b [][]float64, size int) [][]float64 {
	result := make([][]float64, size)
	for i := range result {
		result[i] = make([]float64, size)
	}

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(numCPU)

	rowsPerWorker := size / numCPU
	if rowsPerWorker == 0 {
		rowsPerWorker = 1
	}

	for worker := 0; worker < numCPU; worker++ {
		startRow := worker * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if worker == numCPU-1 {
			endRow = size // Last worker takes any remaining rows
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				for j := 0; j < size; j++ {
					sum := float64(0)
					for k := 0; k < size; k++ {
						sum += a[i][k] * b[k][j]
					}
					result[i][j] = sum
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()
	return result
}

func GraduallyIncreaseMatrixSize() func() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				atomic.AddInt32(&baseMatrixSize, growthRate)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return func() { close(quit) }
}

func PerformComputation(computeDuration prometheus.Histogram, computeOps prometheus.Counter, matrixGenTime prometheus.Histogram, matrixSizeGauge prometheus.Gauge, matrixCount prometheus.Counter) {
	size := getCurrentMatrixSize()
	matrixSizeGauge.Set(float64(size))

	genStart := time.Now()
	matrix1 := GenerateMatrix(size)
	matrix2 := GenerateMatrix(size)
	genDuration := time.Since(genStart).Seconds()
	matrixGenTime.Observe(genDuration)

	start := time.Now()
	MultiplyMatrices(matrix1, matrix2, size)
	duration := time.Since(start).Seconds()

	computeDuration.Observe(duration)
	computeOps.Inc()
	matrixCount.Inc()
}
