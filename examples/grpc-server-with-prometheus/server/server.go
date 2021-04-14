package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	grpc_prometheus "github.com/JST1703/go-grpc-prometheus"
	pb "github.com/JST1703/go-grpc-prometheus/examples/grpc-server-with-prometheus/protobuf"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DemoServiceServer defines a Server.
type DemoServiceServer struct{}

func newDemoServer() *DemoServiceServer {
	return &DemoServiceServer{}
}

// SayHello implements a interface defined by protobuf.
func (s *DemoServiceServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "chimichanga"}, nil
}

// NOTE: Graceful shutdown is missing. Don't use this demo in your production setup.
func main() {
	// Listen an actual port.

	bytemetrics := grpc_prometheus.NewServerByteMetrics()
	reg := prometheus.NewRegistry()
	reg.MustRegister(bytemetrics)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9093))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Create a gRPC Server with gRPC interceptor.
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(bytemetrics.NewServerByteStatsHandler()),
	)

	// Create a new api server.
	demoServer := newDemoServer()

	// Register your service.
	pb.RegisterDemoServiceServer(grpcServer, demoServer)

	// Start your http server for prometheus.
	go func() {
		http.Handle("/bytesize", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		if err := http.ListenAndServe(":9092", nil); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// Start your gRPC server.
	log.Fatal(grpcServer.Serve(lis))
}
