package baseprometheus

import (
	"context"
	"fmt"

	"github.com/newity/glog"
	"github.com/atomyze-foundation/common-component/basemetrics"
)

type createChildMetric[TMetricsBus any] func(childBus, parentBus *TMetricsBus, labels []basemetrics.Label)

type BaseMetricsBus[TMetricsBus any] struct {
	log glog.Logger

	mPrefix        string
	createChildren []createChildMetric[TMetricsBus]
}

func NewBus[TMetricsBus any](ctx context.Context, mPrefix string) *BaseMetricsBus[TMetricsBus] {
	l := glog.FromContext(ctx)
	m := &BaseMetricsBus[TMetricsBus]{
		log:     l,
		mPrefix: mPrefix,
	}

	return m
}

func (m *BaseMetricsBus[TMetricsBus]) AddCounter(
	createChild createChildMetric[TMetricsBus],
	mname, descr string, labels ...basemetrics.LabelName,
) (*Counter, error) {
	c, err := newCounter(m.log, m.getMetricName(mname), descr, labels)
	if err != nil {
		return nil, err
	}
	m.createChildren = append(m.createChildren, createChild)
	return c, nil
}

func (m *BaseMetricsBus[TMetricsBus]) AddGauge(
	createChild createChildMetric[TMetricsBus],
	mname, descr string, labels ...basemetrics.LabelName,
) (*Gauge, error) {
	g, err := newGauge(m.log, m.getMetricName(mname), descr, labels)
	if err != nil {
		return nil, err
	}
	m.createChildren = append(m.createChildren, createChild)
	return g, nil
}

func (m *BaseMetricsBus[TMetricsBus]) AddHisto(
	createChild createChildMetric[TMetricsBus],
	mname, descr string, buckets []float64, labels ...basemetrics.LabelName,
) (*Histo, error) {
	h, err := newHisto(m.log, m.getMetricName(mname), descr, buckets, labels)
	if err != nil {
		return nil, err
	}
	m.createChildren = append(m.createChildren, createChild)
	return h, nil
}

func (m *BaseMetricsBus[TMetricsBus]) getMetricName(mname string) string {
	return fmt.Sprintf("%s%s", m.mPrefix, mname)
}

func (m *BaseMetricsBus[TMetricsBus]) CreateChild(createChildBus func(b *BaseMetricsBus[TMetricsBus]) *TMetricsBus,
	parentBus *TMetricsBus,
	labels ...basemetrics.Label,
) *TMetricsBus {
	baseChildBus := &BaseMetricsBus[TMetricsBus]{
		log:            m.log,
		mPrefix:        m.mPrefix,
		createChildren: m.createChildren,
	}
	childBus := createChildBus(baseChildBus)
	for _, cc := range baseChildBus.createChildren {
		cc(childBus, parentBus, labels)
	}

	return childBus
}
