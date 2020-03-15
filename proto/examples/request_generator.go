package examples

import (
	"math/rand"
	"time"

	matmult "github.com/amwolff/2-lg-mats-bench/gen/go/amwolff/matmult/v1"
)

func GenRequest(cols int, rows int) matmult.MultiplyRequest {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	multiplier := &matmult.Matrix{}
	for c := 0; c < cols; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < rows; r++ {
			column.Coefficients = append(column.Coefficients, r1.Float64())
		}
		multiplier.Columns = append(multiplier.Columns, column)
	}

	multiplicand := &matmult.Matrix{}
	for c := 0; c < rows; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < cols; r++ {
			column.Coefficients = append(column.Coefficients, r1.Float64())
		}
		multiplicand.Columns = append(multiplicand.Columns, column)
	}

	return matmult.MultiplyRequest{
		Multiplier:   multiplier,
		Multiplicand: multiplicand,
	}
}
