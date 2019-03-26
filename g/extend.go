package g

import (
	"context"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

func extendHbr(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (hbrPerf []*MetricValue) {
	counterNameID := CoID()
	var hbsExtend = Extend().Hbr
	extendID := CounterIDByName(counterNameID, hbsExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				hbrPerf = append(hbrPerf, GaugeValue(each.Metric, each.Value))
			}
		}
	}
	return
}

func extendResCPU(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (resPerf []*MetricValue) {
	counterNameID := CoID()
	var resExtend = Extend().Rescpu
	extendID := CounterIDByName(counterNameID, resExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				resPerf = append(resPerf, GaugeValue(each.Metric, each.Value))
			}
		}
	}
	return
}

func extendStoragePath(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (storagePathPerf []*MetricValue) {
	counterNameID := CoID()
	var storagePathExtend = Extend().StoragePath
	extendID := CounterIDByName(counterNameID, storagePathExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags string
				if each.Instance != "" {
					tags = "path=" + each.Instance
				} else {
					tags = ""
				}
				storagePathPerf = append(storagePathPerf, GaugeValue(each.Metric, each.Value, tags))
			}
		}
	}
	return
}

func extendStorageAdapter(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (storageAdapterPerf []*MetricValue) {
	counterNameID := CoID()
	var storageAdapterExtend = Extend().StorageAdapter
	extendID := CounterIDByName(counterNameID, storageAdapterExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags string
				if each.Instance != "" {
					tags = "adapter=" + each.Instance
				} else {
					tags = ""
				}
				storageAdapterPerf = append(storageAdapterPerf, GaugeValue(each.Metric, each.Value, tags))
			}
		}
	}
	return
}

func extendPower(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (powerPerf []*MetricValue) {
	counterNameID := CoID()
	var powerExtend = Extend().Power
	extendID := CounterIDByName(counterNameID, powerExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				powerPerf = append(powerPerf, GaugeValue(each.Metric, each.Value))
			}
		}
	}
	return
}

func extendSys(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (sysPerf []*MetricValue) {
	counterNameID := CoID()
	var sysExtend = Extend().Sys
	extendID := CounterIDByName(counterNameID, sysExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags string
				if each.Instance != "" {
					tags = "dev=" + each.Instance
				} else {
					tags = ""
				}
				sysPerf = append(sysPerf, GaugeValue(each.Metric, each.Value, tags))
			}
		}
	}
	return
}

func extendNet(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (netPerf []*MetricValue) {
	counterNameID := CoID()
	var netExtend = Extend().Net
	extendID := CounterIDByName(counterNameID, netExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags string
				if each.Instance != "" {
					tags = "dev=" + each.Instance
				} else {
					tags = ""
				}
				netPerf = append(netPerf, GaugeValue(each.Metric, each.Value, tags))
			}
		}
	}
	return
}

func extendDisk(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (diskPerf []*MetricValue) {
	counterNameID := CoID()
	var diskExtend = Extend().Disk
	extendID := CounterIDByName(counterNameID, diskExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags string
				if each.Instance != "" {
					tags = "dev=" + each.Instance
				} else {
					tags = ""
				}
				diskPerf = append(diskPerf, GaugeValue(each.Metric, each.Value, tags))
			}
		}
	}
	return
}

func extendCPU(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (cpuPerf []*MetricValue) {
	counterNameID := CoID()
	var cpuExtend = Extend().CPU
	extendID := CounterIDByName(counterNameID, cpuExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags string
				if each.Instance != "" {
					tags = "core=" + each.Instance
				} else {
					tags = ""
				}
				cpuPerf = append(cpuPerf, GaugeValue(each.Metric, each.Value, tags))
			}
		}
	}
	return
}

func extendDatastore(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (datastorePerf []*MetricValue) {
	counterNameID := CoID()
	var datastoreExtend = Extend().Datastore
	extendID := CounterIDByName(counterNameID, datastoreExtend)
	dsWithURL := DsWURL()
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags string
				if each.Instance != "" {
					for _, eachDs := range dsWithURL {
						if strings.Index(eachDs.URL, each.Instance) != -1 {
							tags = "dev=" + eachDs.Datastore
							break
						}
					}
				}
				if each.Instance != "" && tags == "" {
					tags = "dev=" + each.Instance
				}
				datastorePerf = append(datastorePerf, GaugeValue(each.Metric, each.Value, tags))
			}
		}
	}
	return
}

func extendMem(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) (memPerf []*MetricValue) {
	counterNameID := CoID()
	var memExtend = Extend().Mem
	extendID := CounterIDByName(counterNameID, memExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				memPerf = append(memPerf, GaugeValue(each.Metric, each.Value))
			}
		}
	}
	return
}

//EsxiExtendMappers EsxiExtend的Map对象
func EsxiExtendMappers() []EFuncsAndInterval {
	interval := Config().Transfer.Interval
	mappers := []EFuncsAndInterval{
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendHbr,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendResCPU,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendStoragePath,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendStorageAdapter,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendPower,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendSys,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendNet,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendDisk,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendCPU,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendDatastore,
			},
			Interval: interval,
		},
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) []*MetricValue{
				extendMem,
			},
			Interval: interval,
		},
	}
	return mappers
}
