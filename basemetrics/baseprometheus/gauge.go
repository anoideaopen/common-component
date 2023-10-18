package baseprometheus

import (
	"github.com/newity/glog"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/atomyze-foundation/common-component/basemetrics"
)

type Gauge struct {
	*baseMetric
	promMetric *prometheus.GaugeVec
}

func (g *Gauge) Add(v float64, labels ...basemetrics.Label) {
	efLabels := g.mergeLabels(labels)
	m, err := g.promMetric.GetMetricWith(efLabels)
	if err != nil {
		g.log.Errorf("gauge GetMetricWith error: %s", err)
		return
	}
	m.Add(v)
}

func (g *Gauge) Set(v float64, labels ...basemetrics.Label) {
	efLabels := g.mergeLabels(labels)
	m, err := g.promMetric.GetMetricWith(efLabels)
	if err != nil {
		g.log.Errorf("gauge GetMetricWith error: %s", err)
		return
	}
	m.Set(v)
}

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
