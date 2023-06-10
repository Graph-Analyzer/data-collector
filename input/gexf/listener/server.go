package listener

import (
	"context"
	"fmt"
	"graph-analyzer/data-collector/input/gexf/internal"
	"graph-analyzer/data-collector/input/gexf/listener/pb"
	"graph-analyzer/data-collector/repository"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedGexfServiceServer
	pb.UnimplementedHealthCheckServiceServer
	repo repository.GraphRepository
}

func (s *server) ProcessGexf(_ context.Context, in *pb.GexfRequest) (*pb.GexfResponse, error) {
	gexfContent, err := internal.UnmarshalGexf(in.FileContent)
	if err != nil {
		return nil, status.Error(codes.Aborted, "Error unmarshalling GEXF content")
	}

	s.repo.DeleteAll()
	s.repo.CleanupGraph(0)

	internal.GraphWorker(s.repo, gexfContent, in.NetworkName)

	return &pb.GexfResponse{
		Success: true,
	}, nil
}

func (s *server) Check(context.Context, *emptypb.Empty) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Healthy: true,
	}, nil
}

func GrpcServer(repo repository.GraphRepository, host string, port int64) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	// Currently 8 MiB
	s := grpc.NewServer(grpc.MaxRecvMsgSize(8 << 20))
	pb.RegisterGexfServiceServer(s, &server{repo: repo})
	pb.RegisterHealthCheckServiceServer(s, &server{})

	reflection.Register(s)

	log.Debugf("Starting gRPC Server on %s:%d", host, port)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
