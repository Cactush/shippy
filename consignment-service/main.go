package main

import (
	"context"
	"fmt"
	pb "github.com/Cactush/shippy/consignment-service/proto/consignment"
	micro "github.com/micro/go-micro"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo IRepository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}
func main() {
	repo := &Repository{}
	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	srv.Init()
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}

/*
grpc 生成命令
	protoc -I. --go_out=plugins=grpc:.  proto/consignment/consignment.proto
micro 生成命令
    // protoc -I. --go_out=plugins=micro:. proto/consignment/consignment.proto
	protoc --proto_path=. --go_out=. --micro_out=. proto/consignment/consignment.proto
*/
