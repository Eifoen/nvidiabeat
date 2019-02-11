// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period       time.Duration `config:"period"`
	PeriodSystem time.Duration `config:"systemperiod"`
	MetricSets   []string      `config:"metricsets"`
}

var DefaultConfig = Config{
	Period:       1 * time.Second,
	PeriodSystem: 1 * time.Minute,
	MetricSets:   []string{"device", "system"},
}
