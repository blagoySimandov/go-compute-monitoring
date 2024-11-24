package compute

import (
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	MatrixSize = 500 // size of a *square* matrix
)

func GenerateMatrix() [][]float64 {
	matrix := make([][]float64, MatrixSize)
	for i := range matrix {
		matrix[i] = make([]float64, MatrixSize)
		for j := range matrix[i] {
			matrix[i][j] = rand.Float64()
		}
	}
	return matrix
}

func MultiplyMatrices(a, b [][]float64) [][]float64 {
	result := make([][]float64, MatrixSize)
	for i := range result {
		result[i] = make([]float64, MatrixSize)
		for j := range result[i] {
			for k := 0; k < MatrixSize; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

func PerformComputation(computeDuration prometheus.Histogram, computeOps prometheus.Counter, matrixGenTime prometheus.Histogram) {
	genStart := time.Now()
	matrix1 := GenerateMatrix()
	matrix2 := GenerateMatrix()
	genDuration := time.Since(genStart).Seconds()
	matrixGenTime.Observe(genDuration)

	start := time.Now()
	MultiplyMatrices(matrix1, matrix2)
	duration := time.Since(start).Seconds()

	computeDuration.Observe(duration)
	computeOps.Inc()
}
