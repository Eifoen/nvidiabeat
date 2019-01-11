package beater

import (
	"github.com/elastic/beats/libbeat/common"

	"github.com/Eifoen/gonvml"
)

func buildSystemMetrics() (common.MapStr, error) {
	var version string
	var nvml string
	var count uint
	var tmperr, err error

	m := make(common.MapStr)

	// Driver
	version, tmperr = gonvml.SystemDriverVersion()
	if tmperr != nil {
		err = tmperr
	} else {
		m.Put("driver.version", version)
	}

	// NVML
	nvml, tmperr = gonvml.SystemNVMLVersion()
	if tmperr != nil {
		err = tmperr
	} else {
		m.Put("nvml.version", nvml)
	}

	// Device count
	count, tmperr = gonvml.DeviceCount()
	if tmperr != nil {
		err = tmperr
	} else {
		m.Put("device.count", count)
	}

	if len(m) > 0 {
		return common.MapStr{
			"system": m,
		}, err
	}

	return nil, err
}
