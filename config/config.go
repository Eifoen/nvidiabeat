// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period       time.Duration `config:"period"`
	PeriodSystem time.Duration `config:"systemperiod"`
	SystemSet    bool          `config:"system"`
	DeviceSet    bool          `config:"device"`
	ProcessSet   bool          `config:"process"`
}

var DefaultConfig = Config{
	Period:       1 * time.Second,
	PeriodSystem: 1 * time.Minute,
	SystemSet:    true,
	DeviceSet:    true,
	ProcessSet:   false,
}
