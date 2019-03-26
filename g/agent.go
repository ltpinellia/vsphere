package g

import (
	"context"

	"github.com/vmware/govmomi"
)

//AgentMetrics agent存活监控
func AgentMetrics(ctx context.Context, c *govmomi.Client) []*MetricValue {
	return []*MetricValue{GaugeValue("agent.alive", 1)}
}
