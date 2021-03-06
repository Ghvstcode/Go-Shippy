package main

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/micro/go-micro/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/GhvstCode/Go-Shippy/shippy-service-consignment/proto/consignment"
)
const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

//the service struct implements all the methods and is what we pass into the gRPC Registershipping server
type service struct {
	repo repository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	return &pb.Response{Consignments: consignments}, nil
}

func main() {
	repo := &Repository{}

	mService := micro.NewService(
		micro.Name("shippy.service.consignment"),
	)

	// Init will parse the command line flags.
	mService.Init()

	if err := pb.RegisterShippingServiceHandler(mService.Server(), &service{repo}); err != nil {
		log.Panic(err)
	}

	if err := mService.Run(); err != nil {
		log.Panic(err)
	}


	//lis, err := net.Listen("tcp", port)
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//s := grpc.NewServer()
	//
	//pb.RegisterShippingServiceServer(s, &service{repo})
	//reflection.Register(s)
	//
	//log.Println("Running on port:", port)
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
}
