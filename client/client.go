package main

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"math"
	"os"
	matmult "github.com/amwolff/2-lg-mats-bench/gen/go"
	"context"
	"gonum.org/v1/gonum/mat"
	"log"
	"time"
)

const serverAddr = "0.0.0.0:50051"

type goNumMatrices struct {
	Multiplier   mat.Matrix
	Multiplicand mat.Matrix
	Product      mat.Dense
}

func (gn *goNumMatrices) parsePBRequestMessage(message *matmult.PerformerRequest) {
	multiplierRows := len(message.Multiplier.GetColumns()[0].GetCoefficients())
	multiplierColumns := len(message.Multiplier.GetColumns())
	multiplicandColumns := len(message.Multiplicand.GetColumns())

	var v1 []float64
	for i := 0; i < multiplierColumns; i++ {
		for j := 0; j < multiplierRows; j++ {
			v1 = append(v1, message.Multiplier.Columns[j].Coefficients[i])
		}
	}

	var v2 []float64
	for i := 0; i < multiplierRows; i++ {
		for j := 0; j < multiplicandColumns; j++ {
			v2 = append(v2, message.Multiplicand.Columns[j].Coefficients[i])
		}
	}

	gn.Multiplier = mat.NewDense(multiplierRows, multiplierColumns, v1)
	gn.Multiplicand = mat.NewDense(multiplierColumns, multiplicandColumns, v2)
}

func (gn *goNumMatrices) parsePBResponseMessage(message *matmult.PerformerResponse) {
	productRows := len(message.Result.GetColumns()[0].GetCoefficients())
	productColumns := len(message.Result.GetColumns())

	var v []float64
	for i := 0; i < productColumns; i++ {
		for j := 0; j < productRows; j++ {
			v = append(v, message.Result.Columns[j].Coefficients[i])
		}
	}

	gn.Product = *mat.NewDense(productRows, productColumns, v)
}

func getPBResponseMessage(client matmult.PerformerClient, request *matmult.PerformerRequest) (*matmult.PerformerResponse, error) {
	log.Printf("Waiting for response...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.MultiplyMatrices(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func parseJSONFile(path string) (*matmult.PerformerRequest, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var proto matmult.PerformerRequest
	if err := json.Unmarshal(jsonByte, &proto); err != nil {
		return nil, err
	}

	return &proto, nil
}

func matPrint(X ...mat.Matrix) {
	for _, x := range X {
		fa := mat.Formatted(x, mat.Prefix(""), mat.Squeeze())
		fmt.Printf("%v\n", fa)
	}
}

func main() {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt64)))
	if err != nil {
		log.Fatalf("Couldn't connect to the server: %v\n", err)
	}
	defer conn.Close()
	client := matmult.NewPerformerClient(conn)

	request, err := parseJSONFile("/tmp/MACIERZE.JSON")
	if err != nil {
		log.Fatalf("Couldn't parse JSON input file: %v\n", err)
	}

	var goNum goNumMatrices
	goNum.parsePBRequestMessage(request)
	//goNum.Product.Mul(goNum.Multiplier, goNum.Multiplicand)

	response, err := getPBResponseMessage(client, request)
	if err != nil {
		log.Fatalf("Couldn't get a response: %v\n", err)
	}
	goNum.parsePBResponseMessage(response)

	matPrint(goNum.Multiplier, goNum.Multiplicand, &goNum.Product)
}
