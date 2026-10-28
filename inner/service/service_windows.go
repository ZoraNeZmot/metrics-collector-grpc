package service

import (
	"context"
	"errors"
	"fmt"
	"grpc/pkg/metrics"
	"grpc/pkg/proto"
	"io"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateAndStartService(name, address string) error {
	isService, err := svc.IsWindowsService()
	if err != nil {
		return err
	} else if !isService {
		return errors.New("Not a service.")
	}

	service, err := new(name, address)
	if err != nil {
		return err
	}
	if err := svc.Run(name, service); err != nil {
		return err
	}

	return nil
}

func new(name, address string) (*MetricsService, error) {
	client, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	conn := proto.NewMetricsCollectorClient(client)

	eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Info|eventlog.Warning)
	elog, err := eventlog.Open(name)
	if err != nil {
		return nil, err
	}

	return &MetricsService{conn: conn, elog: elog}, nil
}

type MetricsService struct {
	conn proto.MetricsCollectorClient
	elog *eventlog.Log
}

func (ms *MetricsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	accepts := svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}
	// Pending
	changes <- svc.Status{State: svc.Running, Accepts: accepts}
	//start main event loop

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go ms.Stream(ctx)
	for c := range r {
		switch c.Cmd {
		case svc.Stop, svc.Shutdown:
			changes <- svc.Status{State: svc.StopPending}
			return false, 0
		default:
			ms.elog.Error(1, "Wrong MS service command")
		}
	}

	return
}

func (ms *MetricsService) Stream(ctx context.Context) {
	stream, err := ms.conn.StreamMetrics(ctx)
	if err != nil {
		ms.elog.Error(1, "Failed to initiate stream: "+err.Error())
		return
	}
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			m, err := metrics.GetMetrics("192.168.0.1")
			if err != nil {
				return
			}
			err = stream.Send(m)
			if err != nil {
				return
			}
		case <-ctx.Done():
			ack, err := stream.CloseAndRecv()
			if err != nil && !errors.Is(err, io.EOF) {
				ms.elog.Error(1, "Error receiving final Ack: "+err.Error())
			} else if ack != nil {
				ms.elog.Info(1, fmt.Sprintf("Stream closed cleanly. Server Ack success: %v", ack.Success))
			}
			return
		}
	}
}
