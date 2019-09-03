package bench

import (
	"testing"

	matmult "github.com/amwolff/2-lg-mats-bench/gen/go/amwolff/matmult/v1"
	"github.com/golang/protobuf/proto"
)

func BenchmarkPbInPbOut(b *testing.B) {
	conn, client, request, response := setup()
	defer teardown(conn)

	var resp *matmult.MultiplyResponse
	var pbRequest pbRequest
	pbRequest.request = request

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		resp, err = multiplyOnSiteCPP(client, pbRequest)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
	}

	b.StopTimer()
	if !proto.Equal(&response, resp) {
		b.Fatal("Response messages are not equal!")
	}
	b.StartTimer()
}

func BenchmarkPbInGoNumOut(b *testing.B) {
	conn, client, request, response := setup()
	defer teardown(conn)

	var resp *matmult.MultiplyResponse
	var gn goNumMessage
	var pbRequest pbRequest
	pbRequest.request = request

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		resp, err = multiplyOnSiteCPP(client, pbRequest)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		gn.parseResponse(*resp)
	}

	b.StopTimer()
	if !proto.Equal(&response, resp) {
		b.Fatal("Response messages are not equal!")
	}
	b.StartTimer()
}

func BenchmarkPbInGoPrimitiveOut(b *testing.B) {
	conn, client, request, response := setup()
	defer teardown(conn)

	var resp *matmult.MultiplyResponse
	var pm primitiveMessage
	var pbRequest pbRequest
	pbRequest.request = request

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		resp, err = multiplyOnSiteCPP(client, pbRequest)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		pm.parseResponse(*resp)
	}

	b.StopTimer()
	if !proto.Equal(&response, resp) {
		b.Fatal("Response messages are not equal!")
	}
	b.StartTimer()
}

func BenchmarkGoNumInPbOut(b *testing.B) {
	conn, client, request, response := setup()
	defer teardown(conn)

	var resp *matmult.MultiplyResponse
	var gn goNumMessage
	gn.parseRequest(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		resp, err = multiplyOnSiteCPP(client, gn)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
	}

	b.StopTimer()
	if !proto.Equal(&response, resp) {
		b.Fatal("Response messages are not equal!")
	}
	b.StartTimer()
}

func BenchmarkGoNumInGoNumOut(b *testing.B) {
	conn, client, request, response := setup()
	defer teardown(conn)

	var resp *matmult.MultiplyResponse
	var gn goNumMessage
	gn.parseRequest(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		resp, err = multiplyOnSiteCPP(client, gn)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		gn.parseResponse(*resp)
	}

	b.StopTimer()
	if !proto.Equal(&response, resp) {
		b.Fatal("Response messages are not equal!")
	}
	b.StartTimer()
}

func BenchmarkGoNumInGoPrimitiveOut(b *testing.B) {
	conn, client, request, response := setup()
	defer teardown(conn)

	var resp *matmult.MultiplyResponse
	var pm primitiveMessage
	var gn goNumMessage
	gn.parseRequest(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		resp, err = multiplyOnSiteCPP(client, gn)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		pm.parseResponse(*resp)
	}

	b.StopTimer()
	if !proto.Equal(&response, resp) {
		b.Fatal("Response messages are not equal!")
	}
	b.StartTimer()
}

func BenchmarkGoPrimitiveInPbOut(b *testing.B) {
	conn, client, request, response := setup()
	defer teardown(conn)

	var resp *matmult.MultiplyResponse
	var pm primitiveMessage
	pm.parseRequest(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		resp, err = multiplyOnSiteCPP(client, pm)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
	}

	b.StopTimer()
	if !proto.Equal(&response, resp) {
		b.Fatal("Response messages are not equal!")
	}
	b.StartTimer()
}

func BenchmarkGoPrimitiveInGoNumOut(b *testing.B) {
	conn, client, request, response := setup()
	defer teardown(conn)

	var resp *matmult.MultiplyResponse
	var gn goNumMessage
	var pm primitiveMessage
	pm.parseRequest(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		resp, err = multiplyOnSiteCPP(client, pm)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		gn.parseResponse(*resp)
	}

	b.StopTimer()
	if !proto.Equal(&response, resp) {
		b.Fatal("Response messages are not equal!")
	}
	b.StartTimer()
}

func BenchmarkGoPrimitiveInGoPrimitiveOut(b *testing.B) {
	conn, client, request, response := setup()
	defer teardown(conn)

	var resp *matmult.MultiplyResponse
	var pm primitiveMessage
	pm.parseRequest(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		resp, err = multiplyOnSiteCPP(client, pm)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
		pm.parseResponse(*resp)
	}

	b.StopTimer()
	if !proto.Equal(&response, resp) {
		b.Fatal("Response messages are not equal!")
	}
	b.StartTimer()
}

func BenchmarkMultiplyOnSiteGO(b *testing.B) {
	_, _, request, _ := setup()

	var gn goNumMessage
	gn.parseRequest(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gn.Product.Mul(gn.Multiplier, gn.Multiplicand)
	}
}
