package bench

import (
	"context"
	"errors"
	"fmt"
	matmult "github.com/amwolff/2-lg-mats-bench/gen/go/amwolff/matmult/v1"
	"gonum.org/v1/gonum/mat"
	"time"
)

const (
	serverAddr = "0.0.0.0:50051"
	pbPath     = "/tmp/MATRIX"
)

type Writer interface {
	writeRequest() matmult.MultiplyRequest
}

type Parser interface {
	parseRequest(request matmult.MultiplyRequest)
	parseResponse(response matmult.MultiplyResponse)
}

// The data is stored in column-major order
type Matrix struct {
	Data    []float64
	Columns int
	Rows    int
}

type goNumMessage struct {
	Multiplier   mat.Matrix
	Multiplicand mat.Matrix
	Product      mat.Dense
}

type primitiveMessage struct {
	Multiplier   *Matrix
	Multiplicand *Matrix
	Product      *Matrix
}

func (m Matrix) getValueAt(row int, col int) float64 {
	i := col*m.Rows + row
	return m.Data[i]
}

func newMatrix(data []float64, columns int, rows int) *Matrix {
	return &Matrix{data, columns, rows}
}

func (gn goNumMessage) writeRequest() matmult.MultiplyRequest {
	request := &matmult.MultiplyRequest{}
	multiplierRows, multiplierColumns := gn.Multiplier.Dims()
	multiplicandRows, multiplicandColumns := gn.Multiplicand.Dims()

	multiplier := &matmult.Matrix{}
	for c := 0; c < multiplierColumns; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < multiplierRows; r++ {
			column.Coefficients = append(column.Coefficients, gn.Multiplier.At(r, c))
		}
		multiplier.Columns = append(multiplier.Columns, column)
		column = nil
	}

	multiplicand := &matmult.Matrix{}
	for c := 0; c < multiplicandColumns; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < multiplicandRows; r++ {
			column.Coefficients = append(column.Coefficients, gn.Multiplicand.At(r, c))
		}
		multiplicand.Columns = append(multiplicand.Columns, column)
		column = nil
	}

	request.Multiplier = multiplier
	request.Multiplicand = multiplicand

	return *request
}

func (pm primitiveMessage) writeRequest() matmult.MultiplyRequest {
	request := &matmult.MultiplyRequest{}
	multiplierRows, multiplierColumns := pm.Multiplier.Rows, pm.Multiplier.Columns
	multiplicandRows, multiplicandColumns := pm.Multiplicand.Rows, pm.Multiplicand.Columns

	multiplier := &matmult.Matrix{}
	for c := 0; c < multiplierColumns; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < multiplierRows; r++ {
			column.Coefficients = append(column.Coefficients, pm.Multiplier.getValueAt(r, c))
		}
		multiplier.Columns = append(multiplier.Columns, column)
		column = nil
	}

	multiplicand := &matmult.Matrix{}
	for c := 0; c < multiplicandColumns; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < multiplicandRows; r++ {
			column.Coefficients = append(column.Coefficients, pm.Multiplicand.getValueAt(r, c))
		}
		multiplicand.Columns = append(multiplicand.Columns, column)
		column = nil
	}

	request.Multiplier = multiplier
	request.Multiplicand = multiplicand

	return *request
}

func (gn *goNumMessage) parseRequest(request matmult.MultiplyRequest) {
	multiplierRows := len(request.Multiplier.GetColumns()[0].GetCoefficients())
	multiplierColumns := len(request.Multiplier.GetColumns())
	multiplicandRows := len(request.Multiplicand.GetColumns()[0].GetCoefficients())
	multiplicandColumns := len(request.Multiplicand.GetColumns())

	var v1 []float64
	for i := 0; i < multiplierRows; i++ {
		for j := 0; j < multiplierColumns; j++ {
			v1 = append(v1, request.Multiplier.Columns[j].Coefficients[i])
		}
	}

	var v2 []float64
	for i := 0; i < multiplicandRows; i++ {
		for j := 0; j < multiplicandColumns; j++ {
			v2 = append(v2, request.Multiplicand.Columns[j].Coefficients[i])
		}
	}

	gn.Multiplier = mat.NewDense(multiplierRows, multiplierColumns, v1)
	gn.Multiplicand = mat.NewDense(multiplicandRows, multiplicandColumns, v2)
}

func (pm *primitiveMessage) parseRequest(request matmult.MultiplyRequest) {
	multiplierRows := len(request.Multiplier.GetColumns()[0].GetCoefficients())
	multiplierColumns := len(request.Multiplier.GetColumns())
	multiplicandRows := len(request.Multiplicand.GetColumns()[0].GetCoefficients())
	multiplicandColumns := len(request.Multiplicand.GetColumns())

	var v1 []float64
	for i := 0; i < multiplierColumns; i++ {
		for j := 0; j < multiplierRows; j++ {
			v1 = append(v1, request.Multiplier.Columns[i].Coefficients[j])
		}
	}

	var v2 []float64
	for i := 0; i < multiplicandColumns; i++ {
		for j := 0; j < multiplicandRows; j++ {
			v2 = append(v2, request.Multiplicand.Columns[i].Coefficients[j])
		}
	}

	pm.Multiplier = newMatrix(v1, multiplierColumns, multiplierRows)
	pm.Multiplicand = newMatrix(v2, multiplicandColumns, multiplicandRows)
}

func (gn *goNumMessage) parseResponse(response matmult.MultiplyResponse) {
	productRows := len(response.Result.GetColumns()[0].GetCoefficients())
	productColumns := len(response.Result.GetColumns())

	var v []float64
	for i := 0; i < productColumns; i++ {
		for j := 0; j < productRows; j++ {
			v = append(v, response.Result.Columns[j].Coefficients[i])
		}
	}

	gn.Product = *mat.NewDense(productRows, productColumns, v)
}

func (pm *primitiveMessage) parseResponse(response matmult.MultiplyResponse) {
	productRows := len(response.Result.GetColumns()[0].GetCoefficients())
	productColumns := len(response.Result.GetColumns())

	var v []float64
	for i := 0; i < productColumns; i++ {
		for j := 0; j < productRows; j++ {
			v = append(v, response.Result.Columns[j].Coefficients[i])
		}
	}

	pm.Product = newMatrix(v, productColumns, productRows)
}

func multiplyOnSiteCPP(client matmult.MatrixProductAPIClient, request *matmult.MultiplyRequest) (*matmult.MultiplyResponse, error) {
	if len(request.GetMultiplier().GetColumns()) != len(request.GetMultiplicand().GetColumns()) {
		err := errors.New("dimensions mismatch")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.Multiply(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func matPrint(X ...mat.Matrix) {
	for _, x := range X {
		fa := mat.Formatted(x, mat.Prefix(""), mat.Squeeze())
		fmt.Printf("%v\n", fa)
	}
}
