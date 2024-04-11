package baseprometheus

import (
	"github.com/anoideaopen/common-component/basemetrics"
	"github.com/anoideaopen/glog"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// Gauge is a type for gauges
type Gauge struct {
	*baseMetric
	promMetric *prometheus.GaugeVec
}

// Add adds value to gauge
func (g *Gauge) Add(v float64, labels ...basemetrics.Label) {
	efLabels := g.mergeLabels(labels)
	m, err := g.promMetric.GetMetricWith(efLabels)
	if err != nil {
		g.log.Errorf("gauge GetMetricWith error: %s", err)
		return
	}
	m.Add(v)
}

// Set sets value to gauge
func (g *Gauge) Set(v float64, labels ...basemetrics.Label) {
	efLabels := g.mergeLabels(labels)
	m, err := g.promMetric.GetMetricWith(efLabels)
	if err != nil {
		g.log.Errorf("gauge GetMetricWith error: %s", err)
		return
	}
	m.Set(v)
}

// ChildWith creates a child gauge with new labels
func (g *Gauge) ChildWith(labels []basemetrics.Label) *Gauge {
	if bm, ok := g.cloneIfDiffLabels(labels); ok {
		return &Gauge{
			baseMetric: bm,
			promMetric: g.promMetric,
		}
	}
	return g
}

func newGauge(l glog.Logger, name, description string, labels []basemetrics.LabelName) (*Gauge, error) {
	bm := newBaseMetric(l, name, labels)

	promObj := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: description,
	}, bm.labels)

	if err := prometheus.Register(promObj); err != nil {
		return nil, errors.WithStack(err)
	}

	l.Infof("%s gauge created", name)

	return &Gauge{
		promMetric: promObj,
		baseMetric: bm,
	}, nil
}
