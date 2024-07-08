package main

import (
	"grpc-2/pb"
	"time"
)

type server struct {
	pb.UnimplementedProcessesServer
}

func (s *server) GetProcessesInfo(request *pb.ProcessRequest, stream pb.Processes_GetProcessesInfoServer) error {

	processes, err := getProcessInfo()
	if err != nil {
		return err
	}

	for _, process := range processes {
		processName, err := process.Name()
		if err != nil {
			continue
		}

		cpuUsage, err := process.CPUPercent()
		if err != nil {
			continue
		}

		memUsage, err := process.MemoryInfo()
		if err != nil {
			continue
		}

		memUsageInMB := float32(memUsage.RSS) / 1024.0 / 1024.0

		processResponse := &pb.ProcessResponse{
			ProcessId:   uint32(process.Pid),
			ProcessName: processName,
			CpuUsage:    float32(cpuUsage),
			MemoryUsage: memUsageInMB,
			MemUnit:     "MB",
		}
		if err := stream.Send(processResponse); err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 100)
	}
	return nil
}
