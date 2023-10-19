package baseprometheus

import (
	"github.com/atomyze-foundation/common-component/basemetrics"
	"github.com/newity/glog"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Histo struct {
	*baseMetric
	promMetric *prometheus.HistogramVec
}

func (h *Histo) Observe(v float64, labels ...basemetrics.Label) {
	efLabels := h.mergeLabels(labels)
	m, err := h.promMetric.GetMetricWith(efLabels)
	if err != nil {
		h.log.Errorf("histo GetMetricWith error: %s", err)
		return
	}
	m.Observe(v)
}

func (h *Histo) ChildWith(labels []basemetrics.Label) *Histo {
	if bm, ok := h.cloneIfDiffLabels(labels); ok {
		return &Histo{
			baseMetric: bm,
			promMetric: h.promMetric,
		}
	}
	return h
}

func newHisto(l glog.Logger, name, description string, buckets []float64, labels []basemetrics.LabelName) (*Histo, error) {
	bm := newBaseMetric(l, name, labels)

	promObj := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    name,
		Help:    description,
		Buckets: buckets,
	}, bm.labels)

	if err := prometheus.Register(promObj); err != nil {
		return nil, errors.WithStack(err)
	}

	l.Infof("%s histo created", name)

	return &Histo{
		promMetric: promObj,
		baseMetric: bm,
	}, nil
}
