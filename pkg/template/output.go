package template

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"go.k6.io/k6/output"
)

// Output implements the lib.Output interface
type Output struct {
	output.SampleBuffer

	config          Config
	periodicFlusher *output.PeriodicFlusher
	logger          logrus.FieldLogger
}

var _ output.Output = new(Output)

// New creates an instance of the collector
func New(p output.Params) (*Output, error) {
	conf, err := NewConfig(p)
	if err != nil {
		return nil, err
	}
	// Some setupping code

	return &Output{
		config: conf,
		logger: p.Logger,
	}, nil
}

func (o *Output) Description() string {
	return "template: " + o.config.Address
}

func (o *Output) Stop() error {
	o.logger.Debug("Stopping...")
	defer o.logger.Debug("Stopped!")
	o.periodicFlusher.Stop()
	return nil
}

func (o *Output) Start() error {
	o.logger.Debug("Starting...")

	// Here we should connect to a service, open a file or w/e else we decided we need to do

	pf, err := output.NewPeriodicFlusher(o.config.PushInterval, o.flushMetrics)
	if err != nil {
		return err
	}
	o.logger.Debug("Started!")
	o.periodicFlusher = pf

	return nil
}

func (o *Output) flushMetrics() {
	samples := o.GetBufferedSamples()
	start := time.Now()
	var count int
	for _, sc := range samples {
		samples := sc.GetSamples()
		count += len(samples)
		for _, sample := range samples {
			// Here we actually write or accumulate to then write in batches
			// for the template code we just ... dump some parts of it on the screen
			fmt.Printf("%s=%.5f,%s\n", sample.Metric.Name, sample.Value, sample.GetTags().CloneTags())
		}
	}
	if count > 0 {
		o.logger.WithField("t", time.Since(start)).WithField("count", count).Debug("Wrote metrics to stdout")
	}
}
