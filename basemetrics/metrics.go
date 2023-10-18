package basemetrics

type LabelName string

type Label struct {
	Name  LabelName
	Value string
}

func (l LabelName) Create(val string) Label {
	return Label{
		Name:  l,
		Value: val,
	}
}

type Counter interface {
	Inc(labels ...Label)
	Add(v float64, labels ...Label)
}

type Gauge interface {
	Add(v float64, labels ...Label)
	Set(v float64, labels ...Label)
}

type Histogram interface {
	Observe(v float64, labels ...Label)
}
