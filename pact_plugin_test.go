package main

import (
	"context"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	pb "github.com/rob0t7/pact-kafka-plugin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
)

const bufSize = 1024 * 1024

var (
	testListener *bufconn.Listener
	testServer   *grpc.Server
)

func init() {
	// Setup server once for all tests using bufconn for fast in-memory communication
	testListener = bufconn.Listen(bufSize)
	testServer = grpc.NewServer()
	pb.RegisterPactPluginServer(testServer, &pactPluginServer{})
	go func() {
		if err := testServer.Serve(testListener); err != nil {
			// Ignore errors from graceful shutdown
			if err != grpc.ErrServerStopped {
				panic(err)
			}
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return testListener.Dial()
}

func getClient(t *testing.T) (pb.PactPluginClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient("localhost",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client := pb.NewPactPluginClient(conn)
	return client, conn
}

// TestInitPlugin tests the InitPlugin RPC method
func TestInitPlugin(t *testing.T) {
	client, conn := getClient(t)
	// nolint:errcheck
	defer conn.Close()

	tests := []struct {
		name    string
		req     *pb.InitPluginRequest
		wantErr bool
	}{
		{
			name: "valid init request",
			req: &pb.InitPluginRequest{
				Implementation: "pact-jvm",
				Version:        "4.3.0",
			},
			wantErr: false,
		},
		{
			name: "empty implementation",
			req: &pb.InitPluginRequest{
				Implementation: "",
				Version:        "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "empty version",
			req: &pb.InitPluginRequest{
				Implementation: "pact-go",
				Version:        "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.InitPlugin(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitPlugin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp == nil {
				t.Error("expected non-nil response")
			}
			if resp != nil && resp.Catalogue == nil {
				t.Error("expected catalogue to be initialized")
			}

			expected := []*pb.CatalogueEntry{
				{
					Key: PLUGIN_NAME,
					Values: map[string]string{
						"content-types": AVRO_SCHEMA_CONTENT_TYPE,
					},
				},
			}
			if diff := cmp.Diff(expected, resp.Catalogue, protocmp.Transform()); diff != "" {
				t.Errorf("InitPlugin() catalogue mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

// TestUpdateCatalogue tests the UpdateCatalogue RPC method
func TestUpdateCatalogue(t *testing.T) {
	client, conn := getClient(t)
	// nolint:errcheck
	defer conn.Close()

	tests := []struct {
		name    string
		req     *pb.Catalogue
		wantErr bool
	}{
		{
			name: "single entry",
			req: &pb.Catalogue{
				Catalogue: []*pb.CatalogueEntry{
					{
						Type: pb.CatalogueEntry_TRANSPORT,
						Values: map[string]string{
							"content-types": "application/json",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple entries",
			req: &pb.Catalogue{
				Catalogue: []*pb.CatalogueEntry{
					{
						Type: pb.CatalogueEntry_CONTENT_MATCHER,
						Key:  "json",
						Values: map[string]string{
							"content-types": "application/json",
						},
					},
					{
						Type: pb.CatalogueEntry_TRANSPORT,
						Values: map[string]string{
							"content-types": "application/octet-stream",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty catalogue",
			req: &pb.Catalogue{
				Catalogue: []*pb.CatalogueEntry{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.UpdateCatalogue(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCatalogue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp == nil {
				t.Error("expected non-nil google.protobuf.Empty response")
			}
		})
	}
}
