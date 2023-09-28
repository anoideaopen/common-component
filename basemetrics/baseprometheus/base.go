package baseprometheus

import (
	"github.com/atomyze-foundation/common-component/basemetrics"
	"github.com/newity/glog"
)

type baseMetric struct {
	log    glog.Logger
	name   string
	labels []string

	contextLabels map[string]string
}

func newBaseMetric(log glog.Logger, name string, labels []basemetrics.LabelName) *baseMetric {
	bm := &baseMetric{
		log:           log,
		name:          name,
		labels:        labelsNameToStrArr(labels),
		contextLabels: make(map[string]string),
	}

	for _, l := range bm.labels {
		bm.contextLabels[l] = ""
	}
	return bm
}

func (bm *baseMetric) cloneIfDiffLabels(labels []basemetrics.Label) (*baseMetric, bool) {
	contextLabels, isChanged := bm.mergeLabelsExt(labels, false)
	if !isChanged {
		return nil, false
	}
	return &baseMetric{
		log:           bm.log,
		name:          bm.name,
		labels:        bm.labels,
		contextLabels: contextLabels,
	}, true
}

func (bm *baseMetric) mergeLabels(labels []basemetrics.Label) map[string]string {
	res, _ := bm.mergeLabelsExt(labels, true)

	for k, v := range res {
		if v == "" {
			bm.log.Warningf("metric %s label %s with empty value", bm.name, k)
		}
	}
	return res
}

func (bm *baseMetric) mergeLabelsExt(labels []basemetrics.Label, warnIfNotExists bool) (map[string]string, bool) {
	if len(bm.contextLabels) == 0 || len(labels) == 0 {
		return bm.contextLabels, false
	}

	var res map[string]string
	for _, l := range labels {
		if v, ok := bm.contextLabels[string(l.Name)]; !ok || v == l.Value {
			if warnIfNotExists && !ok {
				bm.log.Warningf("metric %s does not contain label %s", bm.name, l.Name)
			}
			continue
		}
		if res == nil {
			res = make(map[string]string)
			for k, v := range bm.contextLabels {
				res[k] = v
			}
		}

		res[string(l.Name)] = l.Value
	}

	if res == nil {
		return bm.contextLabels, false
	}

	return res, true
}

func labelsNameToStrArr(labels []basemetrics.LabelName) []string {
	if len(labels) == 0 {
		return nil
	}
	res := make([]string, 0, len(labels))
	for _, l := range labels {
		res = append(res, string(l))
	}
	return res
}
