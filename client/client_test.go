package main

import (
	matmult "github.com/amwolff/2-lg-mats-bench/gen/go"
	"google.golang.org/grpc"
	"log"
	"math"
	"testing"
)

func BenchmarkParsePBRequestMessage(b *testing.B) {
	request, err := parseJSONFile("/tmp/MACIERZE.JSON")
	if err != nil {
		b.Fatalf("Couldn't parse JSON input file: %v\n", err)
	}

	var goNum goNumMatrices
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		goNum.parsePBRequestMessage(request)
	}
}

func BenchmarkParsePBResponseMessage(b *testing.B) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt64)))
	if err != nil {
		log.Fatalf("Couldn't connect to the server: %v\n", err)
	}
	defer conn.Close()
	client := matmult.NewPerformerClient(conn)
	if err != nil {
		b.Fatalf("Couldn't connect to the server: %v\n", err)
	}

	request, err := parseJSONFile("/tmp/MACIERZE.JSON")
	if err != nil {
		b.Fatalf("Couldn't parse JSON input file: %v\n", err)
	}

	response, err := getPBResponseMessage(client, request)
	if err != nil {
		b.Fatalf("Couldn't get a response: %v\n", err)
	}

	var goNum goNumMatrices
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		goNum.parsePBResponseMessage(response)
	}
}

func BenchmarkParseJSONFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := parseJSONFile("/tmp/MACIERZE.JSON")
		if err != nil {
			b.Errorf("Couldn't parse JSON input file: %v\n", err)
		}
	}
}

func BenchmarkMultiply(b *testing.B) {
	request, err := parseJSONFile("/tmp/MACIERZE.JSON")
	if err != nil {
		b.Fatalf("Couldn't parse JSON input file: %v\n", err)
	}
	var goNum goNumMatrices
	goNum.parsePBRequestMessage(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		goNum.Product.Mul(goNum.Multiplier, goNum.Multiplicand)
	}
}

func BenchmarkGetPBResponseMessage(b *testing.B) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt64)))
	if err != nil {
		log.Fatalf("Couldn't connect to the server: %v\n", err)
	}
	defer conn.Close()
	client := matmult.NewPerformerClient(conn)
	if err != nil {
		b.Fatalf("Couldn't connect to the server: %v\n", err)
	}

	request, err := parseJSONFile("/tmp/MACIERZE.JSON")
	if err != nil {
		b.Fatalf("Couldn't parse JSON input file: %v\n", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = getPBResponseMessage(client, request)
		if err != nil {
			b.Errorf("Couldn't get a response: %v\n", err)
		}
	}
}
