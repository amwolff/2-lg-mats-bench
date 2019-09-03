package bench

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"

	"google.golang.org/grpc"

	matmult "github.com/amwolff/2-lg-mats-bench/gen/go/amwolff/matmult/v1"
	"github.com/golang/protobuf/proto"
	"gonum.org/v1/gonum/mat"
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

type pbRequest struct {
	request matmult.MultiplyRequest
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
	Product      mat.Dense
}

type primitiveMessage struct {
	Multiplier   *matrix
	Multiplicand *matrix
	Product      *matrix
}

func (m matrix) coeff(row int, col int) float64 {
	i := col*m.rows + row
	return m.data[i]
}

func newMatrix(data []float64, columns int, rows int) *matrix {
	return &matrix{data, columns, rows}
}

func (pb pbRequest) writeRequest() matmult.MultiplyRequest {
	return pb.request
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
	multiplierRows, multiplierColumns := pm.Multiplier.rows, pm.Multiplier.columns
	multiplicandRows, multiplicandColumns := pm.Multiplicand.rows, pm.Multiplicand.columns

	multiplier := &matmult.Matrix{}
	for c := 0; c < multiplierColumns; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < multiplierRows; r++ {
			column.Coefficients = append(column.Coefficients, pm.Multiplier.coeff(r, c))
		}
		multiplier.Columns = append(multiplier.Columns, column)
		column = nil
	}

	multiplicand := &matmult.Matrix{}
	for c := 0; c < multiplicandColumns; c++ {
		column := &matmult.Matrix_Column{}
		for r := 0; r < multiplicandRows; r++ {
			column.Coefficients = append(column.Coefficients, pm.Multiplicand.coeff(r, c))
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

func multiplyOnSiteCPP(client matmult.MatrixProductAPIClient, message Writer) (*matmult.MultiplyResponse, error) {
	request := message.writeRequest()
	if len(request.GetMultiplier().GetColumns()) != len(request.GetMultiplicand().GetColumns()) {
		err := errors.New("dimensions mismatch")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.Multiply(ctx, &request)
}

func setup() (*grpc.ClientConn, matmult.MatrixProductAPIClient, matmult.MultiplyRequest, matmult.MultiplyResponse) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt64)))
	if err != nil {
		log.Fatalf("Couldn't connect to the server: %v\n", err)
	}
	client := matmult.NewMatrixProductAPIClient(conn)

	pbFile, err := os.Open(pbPath)
	if err != nil {
		log.Fatalf("Couldn't open protool buffers message: %v\n", err)
	}
	defer pbFile.Close()

	pbBytes, err := ioutil.ReadAll(pbFile)
	if err != nil {
		log.Fatalf("Couldn't read protool buffers message: %v\n", err)
	}

	var pbRequest pbRequest
	var request matmult.MultiplyRequest
	if err := proto.Unmarshal(pbBytes, &request); err != nil {
		log.Fatalf("Couldn't get request message: %v\n", err)
	}
	pbRequest.request = request

	response, err := multiplyOnSiteCPP(client, pbRequest)
	if err != nil {
		log.Fatalf("Couldn't get a response: %v\n", err)
	}

	return conn, client, request, *response
}

func teardown(conn *grpc.ClientConn) {
	conn.Close()
}
