package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc"

	pb "github.com/GhvstCode/Go-Shippy/shippy-service-consignment/proto/consignment"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(fileName string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &consignment); err != nil {
		return nil, err
	}
	return consignment, err
}

func main() {
	//Underthe hood the conn implements ClientConnInterface interface
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewShippingServiceClient(conn)

	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)
}