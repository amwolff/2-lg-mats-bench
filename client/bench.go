package bench

import (
	"context"
	"time"

	matmult "github.com/amwolff/2-lg-mats-bench/gen/go/amwolff/matmult/v1"
	"gonum.org/v1/gonum/mat"
)

type Writer interface {
	writeRequest() *matmult.MultiplyRequest
}

type Parser interface {
	parseRequest(request matmult.MultiplyRequest)
	parseResponse(response matmult.MultiplyResponse)
}

type pbRequest struct {
	request *matmult.MultiplyRequest
}

// The data is stored in column-major order
type matrix struct {
	data    []float64
	columns int
	rows    int
}

type goNumMessage struct {
	Multiplier   mat.Matrix
	Multiplicand mat.Matrix
}

type primitiveMessage struct {
	Multiplier   *matrix
	Multiplicand *matrix
}

func (m matrix) coeff(row int, col int) float64 {
	return m.data[col*m.rows+row]
}

func newMatrix(data []float64, columns int, rows int) *matrix {
	return &matrix{data, columns, rows}
}

func (pb *pbRequest) writeRequest() *matmult.MultiplyRequest {
	return pb.request
}

func (gn *goNumMessage) writeRequest() *matmult.MultiplyRequest {
	return &matmult.MultiplyRequest{
		Multiplier:   fromGonumMatrix(gn.Multiplier),
		Multiplicand: fromGonumMatrix(gn.Multiplicand),
	}
}

func (pm *primitiveMessage) writeRequest() *matmult.MultiplyRequest {
	return &matmult.MultiplyRequest{
		Multiplier:   fromPrimitiveMatrix(pm.Multiplier),
		Multiplicand: fromPrimitiveMatrix(pm.Multiplicand),
	}
}

func (gn *goNumMessage) parseRequest(request *matmult.MultiplyRequest) {
	gn.Multiplier = toGoNum(request.Multiplier)
	gn.Multiplicand = toGoNum(request.Multiplicand)
}

func (pm *primitiveMessage) parseRequest(request *matmult.MultiplyRequest) {
	pm.Multiplier = toPrimitive(request.Multiplier)
	pm.Multiplicand = toPrimitive(request.Multiplicand)
}

func toPrimitive(matrix *matmult.Matrix) *matrix {
	rows := len(matrix.GetColumns()[0].GetCoefficients())
	cols := len(matrix.GetColumns())

	var v []float64
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			v = append(v, matrix.Columns[j].Coefficients[i])
		}
	}

	return newMatrix(v, cols, rows)
}

func toGoNum(matrix *matmult.Matrix) *mat.Dense {
	rows := len(matrix.GetColumns()[0].GetCoefficients())
	cols := len(matrix.GetColumns())

	var v []float64
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			v = append(v, matrix.Columns[j].Coefficients[i])
		}
	}

	return mat.NewDense(rows, cols, v)
}

func fromPrimitiveMatrix(matrix *matrix) *matmult.Matrix {
	retVal := &matmult.Matrix{}
	for c := 0; c < matrix.columns; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < matrix.rows; r++ {
			column.Coefficients = append(column.Coefficients, matrix.coeff(r, c))
		}
		retVal.Columns = append(retVal.Columns, column)
	}

	return retVal
}

func fromGonumMatrix(matrix mat.Matrix) *matmult.Matrix {
	rows, cols := matrix.Dims()

	retVal := &matmult.Matrix{}
	for c := 0; c < cols; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < rows; r++ {
			column.Coefficients = append(column.Coefficients, matrix.At(r, c))
		}
		retVal.Columns = append(retVal.Columns, column)
	}

	return retVal
}

func multiplyAtService(client matmult.MatrixProductAPIClient, message Writer) (*matmult.MultiplyResponse, error) {
	request := message.writeRequest()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return client.Multiply(ctx, request)
}
