package baseprometheus

import (
	"context"

	"github.com/anoideaopen/common-component/basemetrics"
	"github.com/newity/glog"
	"github.com/prometheus/client_golang/prometheus"
)

type ctxKey int

const (
	ctxMetrics ctxKey = iota
)

// NewContext adds Metrics to the Context.
func NewContext(parent context.Context, m *ExampleMetricsBus) context.Context {
	return context.WithValue(parent, ctxMetrics, m)
}

// FromContext gets Logger from the Context.
func FromContext(ctx context.Context) *ExampleMetricsBus {
	if val, ok := ctx.Value(ctxMetrics).(*ExampleMetricsBus); ok {
		return val
	}

	panic("no metrics in context")
}

type LabelNames struct {
	OneLabel   basemetrics.LabelName
	TwoLabel   basemetrics.LabelName
	ThreeLabel basemetrics.LabelName
}

func Labels() LabelNames {
	return allLabels
}

var allLabels = createLabels()

func createLabels() LabelNames {
	return LabelNames{
		OneLabel:   "one_l",
		TwoLabel:   "two_l",
		ThreeLabel: "three_l",
	}
}

type ExampleMetricsBus struct {
	log glog.Logger

	baseBus *BaseMetricsBus[ExampleMetricsBus]

	mCounterOne *Counter
	mGaugeOne   *Gauge
	mHistoOne   *Histo
}

func NewMetrics(ctx context.Context, mPrefix string) (*ExampleMetricsBus, error) {
	l := glog.FromContext(ctx)
	m := &ExampleMetricsBus{
		log:     l,
		baseBus: NewBus[ExampleMetricsBus](ctx, mPrefix),
	}

	var err error

	if m.mCounterOne, err = m.baseBus.AddCounter(
		func(ch, parent *ExampleMetricsBus, labels []basemetrics.Label) {
			ch.mCounterOne = parent.mCounterOne.ChildWith(labels)
		},
		"counter_one", "counter_one descr",
		Labels().OneLabel,
		Labels().TwoLabel); err != nil {
		return nil, err
	}

	if m.mGaugeOne, err = m.baseBus.AddGauge(
		func(ch, parent *ExampleMetricsBus, labels []basemetrics.Label) {
			ch.mGaugeOne = parent.mGaugeOne.ChildWith(labels)
		},
		"gauge_one", "gauge_one descr",
		Labels().TwoLabel, Labels().ThreeLabel); err != nil {
		return nil, err
	}

	if m.mHistoOne, err = m.baseBus.AddHisto(
		func(ch, parent *ExampleMetricsBus, labels []basemetrics.Label) {
			ch.mHistoOne = parent.mHistoOne.ChildWith(labels)
		},
		"histo_one", "histo_one descr",
		prometheus.DefBuckets,
		Labels().ThreeLabel); err != nil {
		return nil, err
	}

	l.Info("prometheus metrics created")
	return m, nil
}

func (m *ExampleMetricsBus) CreateChild(labels ...basemetrics.Label) *ExampleMetricsBus {
	if len(labels) == 0 {
		return m
	}

	return m.baseBus.CreateChild(func(baseChildBus *BaseMetricsBus[ExampleMetricsBus]) *ExampleMetricsBus {
		return &ExampleMetricsBus{
			log:     m.log,
			baseBus: baseChildBus,
		}
	}, m, labels...)
}

func (m *ExampleMetricsBus) CounterOne() basemetrics.Counter { return m.mCounterOne }
func (m *ExampleMetricsBus) GaugeOne() basemetrics.Gauge     { return m.mGaugeOne }
func (m *ExampleMetricsBus) HistoOne() basemetrics.Histogram { return m.mHistoOne }
