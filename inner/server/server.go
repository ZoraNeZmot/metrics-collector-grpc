package server

import (
	pb "grpc/pkg/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type MetricsServer struct {
	port   string
	host   string
	server *grpc.Server

	pb.UnimplementedMetricsCollectorServer
}

func New(port, host string) *MetricsServer {
	return &MetricsServer{
		port: port,
		host: host,
	}

}

func (ms *MetricsServer) Start() error {
	ls, err := net.Listen("tcp", net.JoinHostPort(ms.host, ms.port))
	if err != nil {
		return err
	}
	ms.server = grpc.NewServer()
	pb.RegisterMetricsCollectorServer(ms.server, ms)

	if err := ms.server.Serve(ls); err != nil {
		return err
	}

	return nil
}

func (ms *MetricsServer) StreamMetrics(stream pb.MetricsCollector_StreamMetricsServer) error {
	for {
		payload, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Ack{Success: true})
		} else if err != nil {
			return err
		}

		log.Printf(
			"agent=%s ts=%d cpu=%.2f%% mem=%.2f%% loss=%.2f%%",
			payload.GetAgentId(),
			payload.GetTimestamp(),
			payload.GetCpuUsagePercent(),
			payload.GetMemoryUsagePercent(),
			payload.GetPacketLossPercent(),
		)
	}
}
