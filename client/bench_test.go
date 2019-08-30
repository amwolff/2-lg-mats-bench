package bench

import (
	matmult "github.com/amwolff/2-lg-mats-bench/gen/go/amwolff/matmult/v1"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"math"
	"os"
	"testing"
)

type initVariables struct {
	client   matmult.MatrixProductAPIClient
	request  matmult.MultiplyRequest
	response matmult.MultiplyResponse
}

var initVars initVariables

func TestMain(m *testing.M) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt64)))
	if err != nil {
		log.Fatalf("Couldn't connect to the server: %v\n", err)
	}
	client := matmult.NewMatrixProductAPIClient(conn)

	pbFile, err := os.Open(pbPath)
	if err != nil {
		log.Fatalf("Couldn't open protool buffers message: %v\n", err)
	}

	pbBytes, err := ioutil.ReadAll(pbFile)
	if err != nil {
		log.Fatalf("Couldn't read protool buffers message: %v\n", err)
	}

	var request matmult.MultiplyRequest
	if err := proto.Unmarshal(pbBytes, &request); err != nil {
		log.Fatalf("Couldn't get request message: %v\n", err)
	}

	response, err := multiplyOnSiteCPP(client, &request)
	if err != nil {
		log.Fatalf("Couldn't get a response: %v\n", err)
	}

	initVars.client = client
	initVars.request = request
	initVars.response = *response

	exit := m.Run()
	conn.Close()
	pbFile.Close()
	os.Exit(exit)
}

func BenchmarkPbInPbOut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := multiplyOnSiteCPP(initVars.client, &initVars.request)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
	}
}

func BenchmarkPbInGoNumOut(b *testing.B) {
	var gn goNumMessage

	for i := 0; i < b.N; i++ {
		response, err := multiplyOnSiteCPP(initVars.client, &initVars.request)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		gn.parseResponse(*response)
	}
}

func BenchmarkPbInGoPrimitiveOut(b *testing.B) {
	var pm primitiveMessage

	for i := 0; i < b.N; i++ {
		response, err := multiplyOnSiteCPP(initVars.client, &initVars.request)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		pm.parseResponse(*response)
	}
}

func BenchmarkGoNumInPbOut(b *testing.B) {
	var gn goNumMessage
	gn.parseRequest(initVars.request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := gn.writeRequest()
		_, err := multiplyOnSiteCPP(initVars.client, &req)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
	}
}

func BenchmarkGoNumInGoNumOut(b *testing.B) {
	var gn goNumMessage
	gn.parseRequest(initVars.request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := gn.writeRequest()
		response, err := multiplyOnSiteCPP(initVars.client, &req)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		gn.parseResponse(*response)
	}
}

func BenchmarkGoNumInGoPrimitiveOut(b *testing.B) {
	var pm primitiveMessage
	var gn goNumMessage
	gn.parseRequest(initVars.request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := gn.writeRequest()
		response, err := multiplyOnSiteCPP(initVars.client, &req)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		pm.parseResponse(*response)
	}
}

func BenchmarkGoPrimitiveInPbOut(b *testing.B) {
	var pm primitiveMessage
	pm.parseRequest(initVars.request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := pm.writeRequest()
		_, err := multiplyOnSiteCPP(initVars.client, &req)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
	}
}

func BenchmarkGoPrimitiveInGoNumOut(b *testing.B) {

	var gn goNumMessage
	var pm primitiveMessage
	pm.parseRequest(initVars.request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := pm.writeRequest()
		response, err := multiplyOnSiteCPP(initVars.client, &req)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		gn.parseResponse(*response)
	}
}

func BenchmarkGoPrimitiveInGoPrimitiveOut(b *testing.B) {
	var pm primitiveMessage
	pm.parseRequest(initVars.request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := pm.writeRequest()
		response, err := multiplyOnSiteCPP(initVars.client, &req)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		pm.parseResponse(*response)
	}
}

func BenchmarkMultiplyOnSiteGO(b *testing.B) {
	var gn goNumMessage
	gn.parseRequest(initVars.request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gn.Product.Mul(gn.Multiplier, gn.Multiplicand)
	}
}

func TestParseRequestAndWriteRequest(t *testing.T) {
	var gn goNumMessage
	var pm primitiveMessage

	gn.parseRequest(initVars.request)
	pm.parseRequest(initVars.request)

	gnReq := gn.writeRequest()
	pmReq := pm.writeRequest()

	if !proto.Equal(&initVars.request, &gnReq) {
		t.Error("Given messages are not equal! Something went wrong with processing gonum set of matrices.")
	}

	if !proto.Equal(&initVars.request, &pmReq) {
		t.Error("Given messages are not equal! Something went wrong with processing primitive set of matrices.")
	}
}
