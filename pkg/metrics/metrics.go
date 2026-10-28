package metrics

import (
	pb "grpc/pkg/proto"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

/*
type MetricPayload struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	AgentId            string                 `protobuf:"bytes,1,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	Timestamp          int64                  `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	CpuUsagePercent    float64                `protobuf:"fixed64,3,opt,name=cpu_usage_percent,json=cpuUsagePercent,proto3" json:"cpu_usage_percent,omitempty"`
	MemoryUsagePercent float64                `protobuf:"fixed64,4,opt,name=memory_usage_percent,json=memoryUsagePercent,proto3" json:"memory_usage_percent,omitempty"`
	PacketLossPercent  float64                `protobuf:"fixed64,5,opt,name=packet_loss_percent,json=packetLossPercent,proto3" json:"packet_loss_percent,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}
*/

func GetMetrics(addr string) (*pb.MetricPayload, error) {
	timestamp := time.Now().UnixMilli()
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	networkPercent, err := GetLostPackets(addr)
	if err != nil {
		return nil, err
	}

	return &pb.MetricPayload{
		AgentId:            "RandomStringAsId",
		Timestamp:          timestamp,
		CpuUsagePercent:    cpuPercent[0],
		PacketLossPercent:  networkPercent,
		MemoryUsagePercent: vmStat.UsedPercent,
	}, nil
}

// should be used gateway address
func GetLostPackets(addr string) (float64, error) {
	pinger, err := probing.NewPinger(addr)
	if err != nil {
		return 0, err
	}
	pinger.Count = 10
	pinger.Interval = 100 * time.Millisecond
	pinger.Timeout = 500 * time.Millisecond
	pinger.SetPrivileged(true)

	err = pinger.Run()
	if err != nil {
		return 0, err
	}

	return pinger.Statistics().PacketLoss, nil
}
