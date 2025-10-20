package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	pb "github.com/rob0t7/pact-kafka-plugin/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type pactPluginServer struct {
	pb.UnimplementedPactPluginServer
}

// FindRandomPort finds an available random port by binding to port 0
func FindRandomPort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	// nolint:errcheck
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// StartPluginServer starts the gRPC plugin server on a random port
func StartPluginServer() error {
	// Find a random available port
	port, err := FindRandomPort()
	if err != nil {
		return fmt.Errorf("failed to find random port: %w", err)
	}

	// Generate a random server key
	serverKey := uuid.New().String()

	fmt.Printf(`{"port": %d, "serverKey": "%s"}%s`, port, serverKey, "\n")

	// Start listening on the port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	// Create and start the gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterPactPluginServer(grpcServer, &pactPluginServer{})

	slog.Info("pact-kafka-plugin gRPC server listening", "port", port, "serverKey", serverKey)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *pactPluginServer) InitPlugin(ctx context.Context, req *pb.InitPluginRequest) (*pb.InitPluginResponse, error) {
	slog.Info("Received InitPlugin request:")
	// Return catalogue entries for Kafka plugin
	return &pb.InitPluginResponse{
		Catalogue: []*pb.CatalogueEntry{
			{
				Type: pb.CatalogueEntry_CONTENT_MATCHER,
				Key:  PLUGIN_NAME,
				Values: map[string]string{
					"content-types": AVRO_SCHEMA_CONTENT_TYPE,
				},
			},
		},
	}, nil
}

func (s *pactPluginServer) UpdateCatalogue(ctx context.Context, req *pb.Catalogue) (*emptypb.Empty, error) {
	slog.Info("Received UpdateCatalogue request")
	return &emptypb.Empty{}, nil
}

func (s *pactPluginServer) ConfigureInteraction(ctx context.Context, req *pb.ConfigureInteractionRequest) (*pb.ConfigureInteractionResponse, error) {
	slog.Info("Received ConfigureInteraction request")

	// parse the required fields. The fields should include the schemaID and the message in []byte
	schemaID, ok := req.ContentsConfig.Fields["schemaId"]
	if !ok {
		return nil, fmt.Errorf("schemaId field is required")
	}
	_, ok = schemaID.Kind.(*structpb.Value_NumberValue)
	if !ok {
		return nil, fmt.Errorf("schemaId field must be a number")
	}
	slog.Info("schemaId", "value", int64(schemaID.GetNumberValue()))

	// parse the []byte base64EncodedMessage
	base64EncodedMessage, ok := req.ContentsConfig.Fields["message"]
	if !ok {
		return nil, fmt.Errorf("message field is required")
	}
	_, ok = base64EncodedMessage.Kind.(*structpb.Value_StringValue)
	if !ok {
		return nil, fmt.Errorf("message field must be a string")
	}
	message, err := base64.StdEncoding.DecodeString(base64EncodedMessage.GetStringValue())
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 message: %w", err)
	}

	content := fmt.Sprintf("%X%04X%s", 0, int64(schemaID.GetNumberValue()), message)
	var interactions = make([]*pb.InteractionResponse, 0)
	interaction := &pb.InteractionResponse{
		Contents: &pb.Body{
			ContentType: AVRO_SCHEMA_CONTENT_TYPE,
			Content:     wrapperspb.Bytes([]byte(content)),
		},
		PartName: "message",
	}
	interactions = append(interactions, interaction)

	return &pb.ConfigureInteractionResponse{
		Interaction: interactions,
	}, nil
}

func (s *pactPluginServer) CompareContents(ctx context.Context, req *pb.CompareContentsRequest) (*pb.CompareContentsResponse, error) {
	slog.Info("Received CompareContents request", "request", req)

	expectedContentType := req.Expected.ContentType
	actualContentType := req.Actual.ContentType

	if expectedContentType != actualContentType {
		return &pb.CompareContentsResponse{
			TypeMismatch: &pb.ContentTypeMismatch{
				Actual:   actualContentType,
				Expected: expectedContentType,
			},
		}, nil
	}

	return &pb.CompareContentsResponse{
		Results: map[string]*pb.ContentMismatches{
			"$": {
				Mismatches: []*pb.ContentMismatch{
					{
						Actual:   req.Actual.Content,
						Expected: req.Expected.Content,
						Mismatch: "The content does not match",
					},
				},
			},
		},
	}, nil
}
