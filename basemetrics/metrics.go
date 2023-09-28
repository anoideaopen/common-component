package basemetrics

// LabelName is a type for label names
type LabelName string

// Label is a type for labels
type Label struct {
	Name  LabelName
	Value string
}

// Create creates a new label
func (l LabelName) Create(val string) Label {
	return Label{
		Name:  l,
		Value: val,
	}
}

// Counter is a type for counters
type Counter interface {
	Inc(labels ...Label)
	Add(v float64, labels ...Label)
}

// Gauge is a type for gauges
type Gauge interface {
	Add(v float64, labels ...Label)
	Set(v float64, labels ...Label)
}

// Histogram is a type for histograms
type Histogram interface {
	Observe(v float64, labels ...Label)
}
