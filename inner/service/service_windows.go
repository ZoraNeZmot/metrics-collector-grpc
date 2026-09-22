package service

import (
	"errors"
	pb "grpc/pkg/proto"

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
	conn := pb.NewMetricsCollectorClient(client)

	eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Info|eventlog.Warning)
	elog, err := eventlog.Open(name)
	if err != nil {
		return nil, err
	}

	return &MetricsService{conn: conn, elog: elog}, nil
}

type MetricsService struct {
	conn pb.MetricsCollectorClient
	elog *eventlog.Log
}

func (ms *MetricsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	accepts := svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}
	// Pending
	changes <- svc.Status{State: svc.Running, Accepts: accepts}
	//start main event loop

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
