package beater

import (
	"fmt"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/Eifoen/gonvml"
	"github.com/Eifoen/nvidiabeat/config"
)

const (
	metricsetProcess = "process"
	metricsetDevice  = "device"
	metricsetSystem  = "system"
)

// Nvidiabeat configuration.
type Nvidiabeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of nvidiabeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Nvidiabeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func RunSystem(b *beat.Beat, period time.Duration, channel chan beat.Event, done chan struct{}) error {
	ticker := time.NewTicker(period)

	for {
		select {
		case <-done:
			logp.Info("System: quitting thread")
			return nil
		case <-ticker.C:
			system, err := buildSystemMetrics()
			if err != nil {
				logp.Err("Unable to obtain system related metrics: " + err.Error())
			}
			event := beat.Event{
				Timestamp: time.Now(),
				Fields:    system,
			}
			channel <- event
			logp.Info("System: Event queued")
		}
	}
}

func RunDevice(b *beat.Beat, period time.Duration, channel chan beat.Event, done chan struct{}) error {
	ticker := time.NewTicker(period)

	for {
		select {
		case <-done:
			logp.Info("Device: quitting thread")
			return nil
		case <-ticker.C:
			devices, err := buildDevicesMetrics()
			if err != nil {
				logp.Err("Unable to obtain devices related metrics: " + err.Error())
			}

			for _, device := range devices {
				event := beat.Event{
					Timestamp: time.Now(),
					Fields:    device,
				}
				channel <- event
				logp.Info("Device-Event queued")
			}
		}
	}
}

func RunProcess(b *beat.Beat, period time.Duration, channel chan beat.Event, done chan struct{}) error {
	ticker := time.NewTicker(period)

	for {
		select {
		case <-done:
			logp.Info("Process: quitting thread")
			return nil
		case <-ticker.C:
			//get processes
			procs, err := buildDevicesProcessesMetrics()
			if err != nil {
				logp.Err("Unable to obtain processes related metrics: " + err.Error())
			}

			//bulk send processes
			for _, proc := range procs {
				event := beat.Event{
					Timestamp: time.Now(),
					Fields:    proc,
				}
				channel <- event
				logp.Info("Process-Event sent")
			}
		}
	}
}

func contains(str string, slice []string) bool {
	for _, e := range slice {
		diff := strings.Compare(str, e)
		if diff == 0 {
			return true
		}
	}
	return false
}

// Run starts nvidiabeat.
func (bt *Nvidiabeat) Run(b *beat.Beat) error {
	logp.Info("nvidiabeat is running! Hit CTRL-C to stop it.")
	var err error

	// connect to server
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	// init NVML
	err = gonvml.Initialize()
	if err != nil {
		logp.Critical("Unable to init NVML Library: " + err.Error())
		return fmt.Errorf("Unable to init NVML Library")
	}
	defer gonvml.Shutdown()

	// start workers
	events := make(chan beat.Event, 50)
	quitChannel := make(chan struct{})
	if contains(metricsetSystem, bt.config.MetricSets) {
		go RunSystem(b, bt.config.PeriodSystem, events, quitChannel)
	}
	if contains(metricsetDevice, bt.config.MetricSets) {
		go RunDevice(b, bt.config.Period, events, quitChannel)
	}
	if contains(metricsetProcess, bt.config.MetricSets) {
		go RunProcess(b, bt.config.Period, events, quitChannel)
	}

	logp.Info("Run: started workers")
	// runtime loop
	for {
		select {
		case <-bt.done:
			logp.Info("stopping threads - please wait patiently")
			close(quitChannel)
			return nil
		case event := <-events:
			logp.Info("Run: got from worker event")
			bt.client.Publish(event)
			logp.Info("Run: Event sent")
		}
	}
}

// Stop stops nvidiabeat.
func (bt *Nvidiabeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
