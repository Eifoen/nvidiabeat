package beater

import (
	"fmt"
	"strings"

	"github.com/Eifoen/gonvml"
	"github.com/elastic/beats/libbeat/common"
)

type process struct {
	pid            uint
	graphicsMemory uint64
	computeMemory  uint64
}

type ProcessType int

const (
	ProcessTypeGraphics ProcessType = 1
	ProcessTypeCompute  ProcessType = 2
	ProcessTypeBoth     ProcessType = 3
	ProcessTypeUnknown  ProcessType = -1
)

func (pt ProcessType) String() string {
	switch pt {
	case ProcessTypeBoth:
		return "both"
	case ProcessTypeCompute:
		return "compute"
	case ProcessTypeGraphics:
		return "graphics"
	default:
		return "uknown"
	}
}

func (proc process) Type() (ProcessType, error) {
	if (proc.computeMemory > 0) && (proc.graphicsMemory > 0) {
		return ProcessTypeBoth, nil
	} else if proc.computeMemory > 0 {
		return ProcessTypeCompute, nil
	} else if proc.graphicsMemory > 0 {
		return ProcessTypeGraphics, nil
	}
	return ProcessTypeUnknown, fmt.Errorf("Unable to determine type of process " + string(proc.pid))
}

func buildDevicesProcessesMetrics() ([]common.MapStr, error) {
	var m []common.MapStr
	var err error

	count, tmperr := gonvml.DeviceCount()
	if tmperr != nil {
		return m, tmperr
	}

	for i := uint(0); i < count; i++ {
		device, tmperr := gonvml.DeviceHandleByIndex(i)
		if tmperr == nil {
			metrics, tmperr := buildDeviceProcessesMetrics(device)
			if tmperr == nil {
				for _, metric := range metrics {
					m = append(m, metric)
				}
			}
		} else {
			err = tmperr
		}
	}
	return m, err
}

func buildDeviceProcessesMetrics(dev gonvml.Device) ([]common.MapStr, error) {
	procs := make(map[uint]*process)
	var errors []string

	//compute process
	cprocs, err := dev.ComputeProcesses()
	if err == nil {
		for _, proc := range cprocs {
			procs[proc.PID()] = &process{
				pid:           proc.PID(),
				computeMemory: proc.Memory(),
			}
		}
	} else {
		errors = append(errors, "Unable to query compute processes: "+err.Error())
	}

	//graphics process
	gprocs, err := dev.GraphicsProcesses()
	if err == nil {
		for _, proc := range gprocs {
			mproc, exists := procs[proc.PID()]
			if exists == true {
				mproc.graphicsMemory = proc.Memory()
			} else {
				procs[proc.PID()] = &process{
					pid:            proc.PID(),
					graphicsMemory: proc.Memory(),
				}
			}
		}
	} else {
		errors = append(errors, "Unable to query graphics processes: "+err.Error())
	}

	// generate JSON
	var m []common.MapStr
	for _, proc := range procs {
		uuid, err := dev.UUID()
		if err != nil {
			continue
		}
		tmp := make(common.MapStr)
		tmp.Put("device.uuid", uuid)
		tmp.Put("pid", proc.pid)
		t, err := proc.Type()
		if err == nil {
			tmp.Put("type", t.String())
		}
		tmp.Put("memory.compute.bytes", proc.computeMemory)
		tmp.Put("memory.graphics.bytes", proc.graphicsMemory)
		m = append(m, common.MapStr{
			"process": tmp,
		})
	}

	//Finalze
	if len(errors) > 0 {
		return m, fmt.Errorf(strings.Join(errors, "\n"))
	}
	return m, nil
}
