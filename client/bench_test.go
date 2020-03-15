package bench

import (
	"log"
	"math"
	"testing"

	"2-lg-mats-bench/proto/examples"

	matmult "github.com/amwolff/2-lg-mats-bench/gen/go/amwolff/matmult/v1"
	"gonum.org/v1/gonum/mat"
	"google.golang.org/grpc"
)

const (
	serverAddr   = "0.0.0.0:50051"
	testdataPath = "/tmp/MATRIX"
)

type benchEnv struct {
	conn    *grpc.ClientConn
	client  matmult.MatrixProductAPIClient
	request *matmult.MultiplyRequest
	check   *matmult.MultiplyResponse
}

func (b *benchEnv) teardown() {
	if err := b.conn.Close(); err != nil {
		log.Printf("Close: %v\n", err)
	}
}

func setup(b *testing.B, addr string, testdata string) *benchEnv {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt64)))
	if err != nil {
		b.Fatalf("Couldn't connect to the server: %v\n", err)
	}

	ret := &benchEnv{
		conn:   conn,
		client: matmult.NewMatrixProductAPIClient(conn),
	}

	request := examples.GenRequest(2, 2)
	ret.request = &request

	msg := &goNumMessage{}
	msg.parseRequest(&request)

	r, _ := msg.Multiplier.Dims()
	_, c := msg.Multiplicand.Dims()
	pd := mat.NewDense(r, c, nil)
	pd.Mul(msg.Multiplier, msg.Multiplicand)

	ret.check = &matmult.MultiplyResponse{Result: fromGonumMatrix(pd)}

	return ret
}

var (
	MultiplyResponse  *matmult.MultiplyResponse
	GoNumResponse     *mat.Dense
	PrimitiveResponse *matrix
)

func BenchmarkPbInPbOut(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	b.ResetTimer()
	var response *matmult.MultiplyResponse
	for i := 0; i < b.N; i++ {
		var err error
		resp, err := multiplyAtService(env.client, &pbRequest{request: env.request})
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		//b.StopTimer()
		//if !proto.Equal(env.check, resp) {
		//	b.Fatal("Response messages are not equal!")
		//}
		//b.StartTimer()
		response = resp
	}

	MultiplyResponse = response
}

func BenchmarkPbInGoNumOut(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	b.ResetTimer()
	var response *mat.Dense
	for i := 0; i < b.N; i++ {
		var err error
		resp, err := multiplyAtService(env.client, &pbRequest{request: env.request})
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		//b.StopTimer()
		//if !proto.Equal(env.check, resp) {
		//	b.Fatal("Response messages are not equal!")
		//}
		//b.StartTimer()
		response = toGoNum(resp.Result)
	}

	GoNumResponse = response
}

func BenchmarkPbInGoPrimitiveOut(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	b.ResetTimer()
	var response *matrix
	for i := 0; i < b.N; i++ {
		var err error
		resp, err := multiplyAtService(env.client, &pbRequest{request: env.request})
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		//b.StopTimer()
		//if !proto.Equal(env.check, resp) {
		//	b.Fatal("Response messages are not equal!")
		//}
		//b.StartTimer()
		response = toPrimitive(resp.Result)
	}

	PrimitiveResponse = response
}

func BenchmarkGoNumInPbOut(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	b.ResetTimer()
	var response *matmult.MultiplyResponse
	for i := 0; i < b.N; i++ {
		gn := &goNumMessage{}
		gn.parseRequest(env.request)
		var err error
		resp, err := multiplyAtService(env.client, gn)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		//b.StopTimer()
		//if !proto.Equal(env.check, resp) {
		//	b.Fatal("Response messages are not equal!")
		//}
		//b.StartTimer()
		response = resp
	}

	MultiplyResponse = response
}

func BenchmarkGoNumInGoNumOut(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	b.ResetTimer()
	var response *mat.Dense
	for i := 0; i < b.N; i++ {
		gn := &goNumMessage{}
		gn.parseRequest(env.request)
		var err error
		resp, err := multiplyAtService(env.client, gn)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		//b.StopTimer()
		//if !proto.Equal(env.check, resp) {
		//	b.Fatal("Response messages are not equal!")
		//}
		//b.StartTimer()
		response = toGoNum(resp.Result)
	}

	GoNumResponse = response
}

func BenchmarkGoNumInGoPrimitiveOut(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	b.ResetTimer()
	var response *matrix
	for i := 0; i < b.N; i++ {
		gn := &goNumMessage{}
		gn.parseRequest(env.request)
		var err error
		resp, err := multiplyAtService(env.client, gn)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		//b.StopTimer()
		//if !proto.Equal(env.check, resp) {
		//	b.Fatal("Response messages are not equal!")
		//}
		//b.StartTimer()
		response = toPrimitive(resp.Result)
	}

	PrimitiveResponse = response
}

func BenchmarkGoPrimitiveInPbOut(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	b.ResetTimer()
	var response *matmult.MultiplyResponse
	for i := 0; i < b.N; i++ {
		pm := &primitiveMessage{}
		pm.parseRequest(env.request)
		var err error
		resp, err := multiplyAtService(env.client, pm)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		//b.StopTimer()
		//if !proto.Equal(env.check, resp) {
		//	b.Fatal("Response messages are not equal!")
		//}
		//b.StartTimer()
		response = resp
	}

	MultiplyResponse = response
}

func BenchmarkGoPrimitiveInGoNumOut(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	b.ResetTimer()
	var response *mat.Dense
	for i := 0; i < b.N; i++ {
		pm := &primitiveMessage{}
		pm.parseRequest(env.request)
		var err error
		resp, err := multiplyAtService(env.client, pm)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		//b.StopTimer()
		//if !proto.Equal(env.check, resp) {
		//	b.Fatal("Response messages are not equal!")
		//}
		//b.StartTimer()
		response = toGoNum(resp.Result)
	}

	GoNumResponse = response
}

func BenchmarkGoPrimitiveInGoPrimitiveOut(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	b.ResetTimer()
	var response *matrix
	for i := 0; i < b.N; i++ {
		pm := &primitiveMessage{}
		pm.parseRequest(env.request)
		var err error
		resp, err := multiplyAtService(env.client, pm)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		//b.StopTimer()
		//if !proto.Equal(env.check, resp) {
		//	b.Fatal("Response messages are not equal!")
		//}
		//b.StartTimer()
		response = toPrimitive(resp.Result)
	}

	PrimitiveResponse = response
}

func BenchmarkOnsite(b *testing.B) {
	env := setup(b, serverAddr, testdataPath)
	defer env.teardown()

	gn := &goNumMessage{}
	gn.parseRequest(env.request)
	product := &mat.Dense{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		product.Mul(gn.Multiplier, gn.Multiplicand)
	}
}
