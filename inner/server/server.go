package server

import (
	"grpc/pkg/config"
	pb "grpc/pkg/proto"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

type MetricsServer struct {
	port   string
	host   string
	server *grpc.Server

	files map[string]*os.File
	outs  map[string]*log.Logger
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
	conf := config.Conf
	for {
		payload, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Ack{Success: true})
		} else if err != nil {
			return err
		}

		if conf.IsMultiFiles {

		} else {
			err := ms.printToMainFile(conf, payload)
			if err != nil {
				return err
			}
		}
	}
}

func (ms *MetricsServer) getFlag(conf *config.Config) int {
	flag := os.O_RDWR | os.O_CREATE
	if conf.IsOverride {
		flag |= os.O_TRUNC
	} else {
		flag |= os.O_APPEND
	}
	return flag
}

func (ms *MetricsServer) getOut(conf *config.Config, payload *pb.MetricPayload) (*log.Logger, error) {
	out := ms.outs[payload.AgentId]
	if out == nil {
		flag := ms.getFlag(conf)

		file, err := os.OpenFile(ms.getSpecificFilePath(conf, payload), flag, 0644)
		if err != nil {
			return nil, err
		}

		ms.files[payload.AgentId] = file
		out = log.New(file, payload.AgentId, log.Lmsgprefix)
		ms.outs[payload.AgentId] = out
	}

	return out, nil
}

func (ms *MetricsServer) getMainFilePath(conf *config.Config) string {
	return filepath.Join(conf.DataFolder, conf.OneFileName)
}

func (ms *MetricsServer) getSpecificFilePath(conf *config.Config, payload *pb.MetricPayload) string {
	return filepath.Join(conf.DataFolder, payload.AgentId+".txt")
}

func (ms *MetricsServer) printToMainFile(conf *config.Config, payload *pb.MetricPayload) error {
	out, err := ms.getOut(conf, payload)
	if err != nil {
		return err
	}
	out.Printf("ts=%d cpu=%.2f%% mem=%.2f%% loss=%.2f%%\n",
		payload.Timestamp,
		payload.CpuUsagePercent,
		payload.MemoryUsagePercent,
		payload.PacketLossPercent,
	)

	return nil
}

func (ms *MetricsServer) printToSpecificFile(conf *config.Config, payload *pb.MetricPayload) error {
	out, err := ms.getOut(conf, payload)
	if err != nil {
		return err
	}

	out.Printf("ts=%d cpu=%.2f%% mem=%.2f%% loss=%.2f%%\n",
		payload.Timestamp,
		payload.CpuUsagePercent,
		payload.MemoryUsagePercent,
		payload.PacketLossPercent,
	)
	return nil
}
