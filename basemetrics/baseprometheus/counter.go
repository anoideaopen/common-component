package baseprometheus

import (
	"github.com/anoideaopen/common-component/basemetrics"
	"github.com/newity/glog"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// Counter is a type for counters
type Counter struct {
	*baseMetric
	promMetric *prometheus.CounterVec
}

// Inc increments counter
func (c *Counter) Inc(labels ...basemetrics.Label) {
	efLabels := c.mergeLabels(labels)
	m, err := c.promMetric.GetMetricWith(efLabels)
	if err != nil {
		c.log.Errorf("counter GetMetricWith error: %s", err)
		return
	}
	m.Inc()
}

// Add adds value to counter
func (c *Counter) Add(v float64, labels ...basemetrics.Label) {
	efLabels := c.mergeLabels(labels)
	m, err := c.promMetric.GetMetricWith(efLabels)
	if err != nil {
		c.log.Errorf("counter GetMetricWith error: %s", err)
		return
	}
	m.Add(v)
}

// ChildWith creates a child counter with new labels
func (c *Counter) ChildWith(labels []basemetrics.Label) *Counter {
	if bm, ok := c.cloneIfDiffLabels(labels); ok {
		return &Counter{
			baseMetric: bm,
			promMetric: c.promMetric,
		}
	}
	return c
}

func newCounter(l glog.Logger, name, description string, labels []basemetrics.LabelName) (*Counter, error) {
	bm := newBaseMetric(l, name, labels)

	promObj := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: description,
	}, bm.labels)

	if err := prometheus.Register(promObj); err != nil {
		return nil, errors.WithStack(err)
	}

	l.Infof("%s counter created", name)

	return &Counter{
		promMetric: promObj,
		baseMetric: bm,
	}, nil
}
