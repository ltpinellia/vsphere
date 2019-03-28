package g

import (
	"context"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/vim25/types"
)

// Performance get counter value by counter key
func Performance(ctx context.Context, c *govmomi.Client, MOR types.ManagedObjectReference, counterID int32) ([]*MetricPerf, error) {
	pm := performance.NewManager(c.Client)
	var pQS = perfQuerySpec(MOR, counterID)
	counterKVL, err := pm.Query(ctx, pQS)
	var metricPerf = make([]*MetricPerf, 0)

	if err == nil {
		counterKV, err := pm.ToMetricSeries(ctx, counterKVL)
		if (err == nil) && (counterKV != nil) {
			for _, eachCounter := range counterKV[0].Value {
				metricPerf = append(metricPerf, &MetricPerf{Metric: eachCounter.Name, Value: eachCounter.Value[0], Instance: eachCounter.Instance})
			}
			return metricPerf, nil
		}
	}
	return []*MetricPerf{}, err
}

func perfQuerySpec(MOR types.ManagedObjectReference, counterID int32) []types.PerfQuerySpec {
	var Rqs types.PerfQuerySpec
	Rqs.Entity = MOR
	//var now = time.Now()
	//var now1, _ = time.ParseDuration("-1h")
	//var start = now.Add(now1)
	//Rqs.StartTime = &start
	//Rqs.EndTime = &now
	Rqs.MaxSample = 1
	Rqs.Format = "normal"
	Rqs.IntervalId = 20
	Rqs.MetricId = []types.PerfMetricId{types.PerfMetricId{CounterId: counterID, Instance: "*"}}
	return []types.PerfQuerySpec{Rqs}
}
