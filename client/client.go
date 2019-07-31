package main

import (
	"context"
	"log"
	"time"
	pb "2-lg-mats-bench/server/src/proto"
	//pb "github.com/amwolff/2-lg-mats-bench/server/proto"
	"google.golang.org/grpc"
)

const serverAddr = "0.0.0.0:50051"

func printResult(client pb.PerformerClient, multiplier *pb.Matrix, multiplicand *pb.Matrix) {
	log.Printf("Getting multiplier")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	response, err := client.MultiplyMatrices(ctx, &pb.PerformerRequest{Multiplier: multiplier, Multiplicand: multiplicand})
	if err != nil {
		log.Fatalf("Could not get a result: %v", err)
	}
	log.Println(response.Result)
}


func main() {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to the server: %v", err)
	}
	defer conn.Close()
	client := pb.NewPerformerClient(conn)

	printResult(client, &pb.Matrix{}, &pb.Matrix{})
}
