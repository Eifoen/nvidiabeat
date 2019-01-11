package beater

import (
	"fmt"
	"strings"

	"github.com/Eifoen/gonvml"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

func buildDevicesMetrics() ([]common.MapStr, error) {
	var m []common.MapStr
	var err error

	count, tmperr := gonvml.DeviceCount()
	if tmperr != nil {
		return m, tmperr
	}

	for i := uint(0); i < count; i++ {
		device, tmperr := gonvml.DeviceHandleByIndex(i)
		if tmperr == nil {
			metrics, tmperr := buildDeviceMetrics(device)
			if tmperr == nil {
				m = append(m, metrics)
			}
		} else {
			err = tmperr
		}
	}
	return m, err
}

func buildDeviceMetrics(dev gonvml.Device) (common.MapStr, error) {
	m := make(common.MapStr)

	//ID
	id, err := dev.Index()
	if err == nil {
		m.Put("index", id)
	} else {
		logp.Err("Unable to query device Index: " + err.Error())
	}

	//ProductBrand
	brand, err := dev.Brand()
	if err == nil {
		m.Put("brand", brand.String())
	} else {
		logp.Err("Unable to query device brand: " + err.Error())
	}

	//Product Name
	name, err := dev.Name()
	if err == nil {
		m.Put("name", name)
	} else {
		logp.Err("Unable to query device name: " + err.Error())
	}

	//UUID
	uuid, err := dev.UUID()
	if err == nil {
		m.Put("uuid", uuid)
	} else {
		logp.Err("Unable to query device UUID: " + err.Error())
	}

	//BoardID
	boardid, err := dev.BoardID()
	if err == nil {
		m.Put("board.id", boardid)
	} else {
		logp.Err("Unable to query device board id: " + err.Error())
	}

	//ComputeMode
	cm, err := dev.ComputeMode()
	if err == nil {
		m.Put("computemode", cm.String())
	} else {
		logp.Err("Unable to query device compute mode: " + err.Error())
	}

	//FanSpeed
	speed, err := dev.FanSpeed()
	if err == nil {
		m.Put("fanspeed.pct", (float64(speed) / 100))
	} else {
		logp.Err("Unable to query device fan speed: " + err.Error())
	}

	//VBiosVersion
	vbios, err := dev.VBiosVersion()
	if err == nil {
		m.Put("vbios.version", vbios)
	} else {
		logp.Err("Unable to query device VBios version: " + err.Error())
	}

	//Serial
	serial, err := dev.Serial()
	if err == nil {
		m.Put("serial", serial)
	} else {
		logp.Err("Unable to query device serial: " + err.Error())
	}

	//PCIe
	pcie, err := buildDevicePCIeMetrics(dev)
	if err == nil {
		m.Put("pcie", pcie)
	} else {
		logp.Err("Unable to query device pcie interface: " + err.Error())
	}

	//Temperature
	temp, err := buildDeviceTemperatureMetrics(dev)
	if err == nil {
		m.Put("temperature", temp)
	} else {
		logp.Err("Unable to query device temperature: " + err.Error())
	}

	//Power
	power, err := buildDevicePowerMetrics(dev)
	if err == nil {
		m.Put("power", power)
	} else {
		logp.Err("Unable to query device power: " + err.Error())
	}

	//Display
	display, err := buildDeviceDisplayMetrics(dev)
	if err == nil {
		m.Put("display", display)
	} else {
		logp.Err("Unable to query device display metrics: " + err.Error())
	}

	//Memory
	memory, err := buildDeviceMemoryMetrics(dev)
	if err == nil {
		m.Put("memory", memory)
	} else {
		logp.Err("Unable to query device memory metrics: " + err.Error())
	}

	//Utilization
	utilization, err := buildDeviceUtilMetrics(dev)
	if err == nil {
		m.Put("utilization", utilization)
	} else {
		logp.Err("Unable to query device utilization metrics: " + err.Error())
	}

	//Clock
	clocks, err := buildDeviceClocks(dev)
	if err == nil {
		m.Put("clock", clocks)
	} else {
		logp.Err("Unable to query device clocks: " + err.Error())
	}

	//Test
	t := make(common.MapStr)
	t.CopyFieldsTo(m, "test")

	return common.MapStr{
		"device": m,
	}, nil
}

func buildDeviceMemoryMetrics(dev gonvml.Device) (common.MapStr, error) {
	m := make(common.MapStr)
	var errors []string

	//total ande used bytes
	total, used, err := dev.MemoryInfo()
	if err == nil {
		free := total - used

		//percentages
		pctUsed := float64(used) / float64(total)
		pctFree := float64(free) / float64(total)

		m.Put("total.bytes", total)

		m.Put("free.bytes", free)
		m.Put("free.pct", pctFree)

		m.Put("used.bytes", used)
		m.Put("used.pct", pctUsed)

	} else {
		errors = append(errors, "Unable to query device memory "+err.Error())
	}

	//BAR1
	total, used, err = dev.Bar1MemoryInfo()
	if err == nil {
		free := total - used

		//percentages
		pctUsed := float64(used) / float64(total)
		pctFree := float64(free) / float64(total)

		bar1 := make(common.MapStr)

		bar1.Put("total.bytes", total)

		bar1.Put("free.bytes", free)
		bar1.Put("free.pct", pctFree)

		bar1.Put("used.bytes", used)
		bar1.Put("used.pct", pctUsed)

		m.Put("bar1", bar1)
	} else {
		errors = append(errors, "Unable to query device BAR1 memory"+err.Error())
	}
	if len(errors) > 0 {
		return m, fmt.Errorf(strings.Join(errors, "\n"))
	}
	return m, nil
}

func buildDeviceUtilMetrics(dev gonvml.Device) (common.MapStr, error) {
	m := make(common.MapStr)
	var err error

	//GPU & Memory
	gpuUtil, memUtil, tmperr := dev.UtilizationRates()
	if err == nil {
		m.Put("gpu.pct", (float32(gpuUtil) / 100))
		m.Put("memory.pct", (float32(memUtil) / 100))
	} else {
		err = tmperr
	}

	//encoder
	encoderUtil, _, tmperr := dev.EncoderUtilization()
	if err == nil {
		m.Put("encoder.pct", (float32(encoderUtil) / 100))
	} else {
		err = tmperr
	}

	//decoder
	decoderUtil, _, tmperr := dev.DecoderUtilization()
	if err == nil {
		m.Put("decoder.pct", (float32(decoderUtil) / 100))
	} else {
		err = tmperr
	}
	return m, err
}

func buildDeviceDisplayMetrics(dev gonvml.Device) (common.MapStr, error) {
	m := make(common.MapStr)
	var err error

	//DisplayMode
	dm, tmperr := dev.DisplayMode()
	if err == nil {
		m.Put("mode", dm.String())
	} else {
		err = tmperr
	}

	//DisplayActive
	da, tmperr := dev.DisplayActive()
	if err == nil {
		m.Put("active", da.String())
	} else {
		err = tmperr
	}

	return m, err
}

func buildDeviceTemperatureMetrics(dev gonvml.Device) (common.MapStr, error) {
	m := make(common.MapStr)
	var errors []string

	//value
	value, err := dev.Temperature()
	if err == nil {
		m.Put("value", value)
	} else {
		errors = append(errors, "Unable to query device temperature: "+err.Error())
	}

	//Thresholds
	shutdown, slowdown, err := dev.TemperatureThresholds()
	if err == nil {
		m.Put("threshold.shutdown", shutdown)
		m.Put("threshold.slowdown", slowdown)
	} else {
		errors = append(errors, "Unable to query device temperature thresholds: "+err.Error())
	}

	if len(errors) > 0 {
		return m, fmt.Errorf(strings.Join(errors, "\n"))
	}
	return m, nil
}

func buildDevicePowerMetrics(dev gonvml.Device) (common.MapStr, error) {
	m := make(common.MapStr)
	var errors []string

	// PowerState
	ps, err := dev.PowerState()
	if err == nil {
		m.Put("state", ps.String())
	} else {
		errors = append(errors, "Unable to query power state: "+err.Error())
	}

	// PowerUsage
	usage, err := dev.PowerUsage()
	if err == nil {
		m.Put("usage.mW", usage)
	} else {
		errors = append(errors, "Unable to query power usage: "+err.Error())
	}

	//Limits
	management, enforced, err := dev.PowerLimits()
	if err == nil {
		m.Put("limit.management.mW", management)
		m.Put("limit.enforced.mW", enforced)
	} else {
		errors = append(errors, "Unable to query power limits: "+err.Error())
	}

	// Finalize
	if len(errors) > 0 {
		return m, fmt.Errorf(strings.Join(errors, "\n"))
	}
	return m, nil
}

func buildDevicePCIeMetrics(dev gonvml.Device) (common.MapStr, error) {
	m := make(common.MapStr)
	var errors []string

	//Utilization
	util, err := buildDevicePCIeUtilization(dev)
	if err == nil {
		m.Put("utilization", util)
	} else {
		errors = append(errors, err.Error())
	}

	//Link
	link, err := buildDevicePCIeLink(dev)
	if err == nil {
		m.Put("link", link)
	} else {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return m, fmt.Errorf(strings.Join(errors, "\n"))
	}
	return m, nil
}

func buildDevicePCIeUtilization(dev gonvml.Device) (common.MapStr, error) {
	m := make(common.MapStr)

	tx, rx, err := dev.PCIeThroughput()
	if err == nil {
		m.Put("tx.kb", tx)
		m.Put("rx.kb", rx)
		return m, nil
	}
	return m, err
}

func buildDevicePCIeLink(dev gonvml.Device) (common.MapStr, error) {
	m := make(common.MapStr)
	var errors []string

	// Link Generation
	curr, max, err := dev.PCIeLinkGen()
	if err == nil {
		m.Put("generation.value", curr)
		m.Put("generation.max", max)
	} else {
		errors = append(errors, "Unable to query PCIe link generation: "+err.Error())
	}

	// Link Generation
	curr, max, err = dev.PCIeLinkWidth()
	if err == nil {
		m.Put("width.value", curr)
		m.Put("width.max", max)
	} else {
		errors = append(errors, "Unable to query PCIe link width: "+err.Error())
	}

	// Finalize
	if len(errors) > 0 {
		return m, fmt.Errorf(strings.Join(errors, "\n"))
	}
	return m, nil
}

func buildDeviceClocks(dev gonvml.Device) (common.MapStr, error) {
	var errors []string
	m := make(common.MapStr)
	cts := []gonvml.ClockType{gonvml.ClockTypeGraphics, gonvml.ClockTypeSM, gonvml.ClockTypeMem}
	for _, ct := range cts {

		clock, err := dev.Clock(ct)
		if err == nil {
			m.Put(ct.String(), clock)
		} else {
			errors = append(errors, "Unable to query "+ct.String()+" clock: "+err.Error())
		}
	}

	// Finalize
	if len(errors) > 0 {
		return m, fmt.Errorf(strings.Join(errors, "\n"))
	}
	return m, nil
}
