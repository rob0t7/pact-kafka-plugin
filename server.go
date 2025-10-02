package main

import (
	"context"
	"log/slog"

	pb "github.com/rob0t7/pact-kafka-plugin/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type pactPluginServer struct {
	pb.UnimplementedPactPluginServer
}

func (s *pactPluginServer) InitPlugin(ctx context.Context, req *pb.InitPluginRequest) (*pb.InitPluginResponse, error) {
	slog.Info("Received InitPlugin request:")
	// Return catalogue entries for Kafka plugin
	return &pb.InitPluginResponse{
		Catalogue: []*pb.CatalogueEntry{
			{
				Type: pb.CatalogueEntry_CONTENT_GENERATOR,
				Key:  "kafka",
				Values: map[string]string{
					"content-types": "application/vnd.kafka.avro.v1+json",
				},
			},
		},
	}, nil
}

func (s *pactPluginServer) UpdateCatalogue(ctx context.Context, req *pb.Catalogue) (*emptypb.Empty, error) {
	slog.Info("Received UpdateCatalogue request")
	return &emptypb.Empty{}, nil
}
