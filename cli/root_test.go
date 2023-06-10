package cli

import (
	"context"
	"fmt"
	"graph-analyzer/data-collector/input/gexf/listener/pb"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_Commands(t *testing.T) {
	ctx := context.Background()

	neo4j, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "neo4j:4.4.19-community",
			ExposedPorts: []string{"7687/tcp"},
			Env: map[string]string{
				"NEO4J_apoc_export_file_enabled":              "true",
				"NEO4J_apoc_import_file_enabled":              "true",
				"NEO4J_dbms_security_procedures_unrestricted": "apoc.*,algo.*,gds.*",
				"NEO4J_apoc_uuid_enabled":                     "true",
				"NEO4J_dbms_default__listen__address":         "0.0.0.0",
				"NEO4J_dbms_allow__upgrade":                   "true",
				"NEO4J_dbms_default__database":                "neo4j",
				"NEO4J_AUTH":                                  "neo4j/test4sa!",
				"NEO4JLABS_PLUGINS":                           "[\"apoc\", \"graph-data-science\"]",
			},
			WaitingFor: wait.ForLog("Bolt enabled"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer func(neo4j testcontainers.Container, ctx context.Context) {
		err := neo4j.Terminate(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(neo4j, ctx)

	currentDir, _ := os.Getwd()

	camouflage, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "shubhendumadhukar/camouflage:0.12.0",
			ExposedPorts: []string{"4312/tcp"},
			Mounts: []testcontainers.ContainerMount{
				testcontainers.BindMount(fmt.Sprintf("%s/../camouflage", currentDir), "/app"),
			},
			WaitingFor: wait.ForLog("Worker sharing HTTP server at"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	neo4jIp, err := neo4j.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}
	neo4jPort, err := neo4j.MappedPort(ctx, "7687")
	if err != nil {
		log.Fatal(err)
	}

	camouflageIp, err := camouflage.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}
	camouflagePort, err := camouflage.MappedPort(ctx, "4312")
	if err != nil {
		log.Fatal(err)
	}

	// GEXF

	rootCmd.SetArgs([]string{
		"gexf",
		"--neo4j-host",
		fmt.Sprintf("bolt://%s", neo4jIp),
		"--neo4j-port",
		neo4jPort.Port(),
		"--neo4j-user",
		"neo4j",
		"--neo4j-password",
		"test4sa!",
		"--file",
		"../testgraph/graph_sa.gexf",
	})

	assert.NotPanics(t, func() {
		err := rootCmd.Execute()
		if err != nil {
			log.Fatal(err)
		}
	})

	// Reset the flag for the next command
	// https://github.com/spf13/cobra/blob/4590150168e93f4b017c6e33469e26590ba839df/completions_test.go#L185
	file := gexfCommand.Flags().Lookup("file")
	file.Changed = false

	rootCmd.SetArgs([]string{
		"gexf",
		"--neo4j-host",
		fmt.Sprintf("bolt://%s", neo4jIp),
		"--neo4j-port",
		neo4jPort.Port(),
		"--neo4j-user",
		"neo4j",
		"--neo4j-password",
		"test4sa!",
		"--listener",
	})

	// As this command won't terminate it is run in a separate goroutine
	go func() {
		assert.NotPanics(t, func() {
			err := rootCmd.Execute()
			if err != nil {
				log.Fatal(err)
			}
		})
	}()
	time.Sleep(10 * time.Second)
	conn, err := grpc.DialContext(context.Background(), "127.0.0.1:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	if err != nil {
		t.Fatalf("Failed to dial gRPC: %v", err)
	}
	defer conn.Close()

	// Create a client and send a health check request.
	healthClient := pb.NewHealthCheckServiceClient(conn)
	resp, err := healthClient.Check(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	// Check the health check response against the expected response.
	assert.Equal(t, pb.HealthCheckResponse{Healthy: true}.Healthy, resp.Healthy)

	// Create a client and send a GEXF request.
	validGexf := strings.NewReader(`
		<?xml version='1.0' encoding='utf-8'?>
		<gexf xmlns="http://www.gexf.net/1.2draft" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.gexf.net/1.2draft http://www.gexf.net/1.2draft/gexf.xsd" version="1.2">
		  <meta lastmodifieddate="2022-10-06">
			<creator>NetworkX 2.8.7</creator>
		  </meta>
		  <graph defaultedgetype="undirected" mode="static" name="">
			<nodes>
			  <node id="1" label="XR-1">
			  </node>
			  <node id="2" label="XR-2">
			  </node>
			  <node id="3" label="XR-3">
			  </node>
			  <node id="4" label="XR-4">
			  </node>
			  <node id="5" label="XR-5">
			  </node>
			  <node id="6" label="XR-6">
			  </node>
			  <node id="7" label="XR-7">
			  </node>
			  <node id="8" label="XR-8">
			  </node>
			</nodes>
			<edges>
			  <edge source="1" target="3" id="0" weight="1">
			  </edge>
			  <edge source="1" target="2" id="1" weight="1">
			  </edge>
			  <edge source="2" target="3" id="2" weight="1">
			  </edge>
			  <edge source="2" target="4" id="3" weight="1">
			  </edge>
			  <edge source="4" target="3" id="5" weight="1">
			  </edge>
			  <edge source="4" target="5" id="6" weight="1">
			  </edge>
			  <edge source="6" target="7" id="7" weight="1">
			  </edge>
			  <edge source="6" target="8" id="8" weight="1">
			  </edge>
			  <edge source="7" target="8" id="9" weight="1">
			  </edge>
			  <edge source="2" target="6" id="10" weight="1">
			  </edge>
			  <edge source="3" target="7" id="11" weight="1">
			  </edge>
			  <edge source="5" target="6" id="12" weight="1">
			  </edge>
			  <edge source="5" target="7" id="13" weight="1">
			  </edge>
			</edges>
		  </graph>
		</gexf>
	`)
	data, err := io.ReadAll(validGexf)
	if err != nil {
		t.Fatalf("Error reading gexf: %v", err)
	}

	gexfClient := pb.NewGexfServiceClient(conn)
	gexfResp, err := gexfClient.ProcessGexf(context.Background(), &pb.GexfRequest{
		FileContent: data,
		NetworkName: "Test",
	})
	if err != nil {
		t.Fatalf("GEXF upload failed: %v", err)
	}

	assert.Equal(t, pb.GexfResponse{Success: true}.Success, gexfResp.Success)

	// Test for errors
	_, err = gexfClient.ProcessGexf(context.Background(), &pb.GexfRequest{
		FileContent: []byte("<xml>Invalid GEXF</xml>"),
		NetworkName: "Test",
	})

	assert.Error(t, err)

	statusCode, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.Aborted, statusCode.Code())
	assert.Equal(t, "Error unmarshalling GEXF content", statusCode.Message())

	// Jalapeno

	rootCmd.SetArgs([]string{
		"jalapeno",
		"--neo4j-host",
		fmt.Sprintf("bolt://%s", neo4jIp),
		"--neo4j-port",
		neo4jPort.Port(),
		"--neo4j-user",
		"neo4j",
		"--neo4j-password",
		"test4sa!",
		"--jagw-host",
		camouflageIp,
		"--jagw-request-port",
		camouflagePort.Port(),
		"--jagw-subscription-port",
		camouflagePort.Port(),
	})

	// As this command won't terminate it is run in a separate goroutine
	go func() {
		assert.NotPanics(t, func() {
			err := rootCmd.Execute()
			if err != nil {
				log.Fatal(err)
			}
		})
	}()
	time.Sleep(30 * time.Second)
}
