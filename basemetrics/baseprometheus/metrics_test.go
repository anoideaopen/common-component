package baseprometheus

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/anoideaopen/common-component/testshlp"
	"github.com/stretchr/testify/require"
)

func TestCreateMetrics(t *testing.T) {
	m, err := NewMetrics(context.Background(), "TestCreateMetrics_")
	require.NoError(t, err)
	require.NotNil(t, m)
	checkMembers(t, m)
}

func TestCreateChildWithLabels(t *testing.T) {
	m, err := NewMetrics(context.Background(), "TestCreateChildWithLabels_")
	require.NoError(t, err)
	require.NotNil(t, m)

	m2 := m.CreateChild(
		Labels().OneLabel.Create("xxx"),
	)
	require.NotNil(t, m2)
	checkMembers(t, m2)
}

func TestContextPutAndGet(t *testing.T) {
	ctx, _ := testshlp.CreateCtxLogger(t)

	m, err := NewMetrics(ctx, "TestContextPutAndGet_")
	require.NoError(t, err)
	require.NotNil(t, m)

	ctx = NewContext(ctx, m)

	m2 := FromContext(ctx)

	require.NotNil(t, m2)
}

func TestContextCounterInc(t *testing.T) {
	ctx, _ := testshlp.CreateCtxLogger(t)

	m, err := NewMetrics(ctx, "TestContextCounterInc_")
	require.NoError(t, err)
	require.NotNil(t, m)

	m.CounterOne().Inc(
		Labels().OneLabel.Create("xxx"),
		Labels().TwoLabel.Create("yyy"),
		Labels().ThreeLabel.Create("true"),
	)
}

func TestLabelsAndChild(t *testing.T) {
	ctx, _ := testshlp.CreateCtxLogger(t)
	m, err := NewMetrics(ctx, "TestLabelsAndChild_")
	require.NoError(t, err)
	require.NotNil(t, m)

	m.CounterOne().Inc(
		Labels().OneLabel.Create("m"),
		Labels().TwoLabel.Create("true"),
	)

	m2 := m.CreateChild(
		Labels().OneLabel.Create("ch-m2"),
	)

	m2.CounterOne().Add(10,
		Labels().TwoLabel.Create("true"),
	)

	m2.CounterOne().Add(10,
		Labels().OneLabel.Create("m2"),
		Labels().TwoLabel.Create("true"),
	)

	m3 := m2.CreateChild(
		Labels().OneLabel.Create("ch-m3"),
	)

	m2.CounterOne().Add(11,
		Labels().TwoLabel.Create("true"),
	)

	m.CounterOne().Inc(
		Labels().OneLabel.Create("m"),
		Labels().TwoLabel.Create("true"),
	)

	m3.CounterOne().Inc(
		Labels().TwoLabel.Create("true"),
	)

	m2.CounterOne().Add(100,
		Labels().OneLabel.Create(""),
		Labels().TwoLabel.Create("true"),
	)
}

// TestMetricsInParallels need to run with -race flag
func TestMetricsInParallels(t *testing.T) {
	ctx, _ := testshlp.CreateCtxLogger(t)
	m, err := NewMetrics(ctx, "TestMetricsInParallels_")
	require.NoError(t, err)
	require.NotNil(t, m)

	mroot := m.CreateChild(
		Labels().OneLabel.Create("root"),
	)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			mroot.CounterOne().Inc(
				Labels().TwoLabel.Create("true"),
			)
		}
	}()

	ctx = NewContext(ctx, mroot)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(ctx context.Context, i int) {
			defer wg.Done()
			mch := FromContext(ctx).CreateChild(
				Labels().OneLabel.Create(fmt.Sprintf("ch-%v", i)),
			)
			for i := 0; i < 1000; i++ {
				mch.CounterOne().Inc(
					Labels().TwoLabel.Create("true"),
				)
			}
		}(ctx, i)
	}

	wg.Wait()
}

func checkMembers(t *testing.T, m *ExampleMetricsBus) {
	require.NotNil(t, m.CounterOne())
	require.Equal(t, m.CounterOne(), m.CounterOne())

	require.NotNil(t, m.GaugeOne())
	require.Equal(t, m.GaugeOne(), m.GaugeOne())

	require.NotNil(t, m.HistoOne())
	require.Equal(t, m.HistoOne(), m.HistoOne())
}
