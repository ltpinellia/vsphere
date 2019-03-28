package g

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

//EsxiList 用于获取VSphere下所有的esxi
func EsxiList(ctx context.Context, c *govmomi.Client) []mo.HostSystem {
	m := view.NewManager(c.Client)
	var esxiList []mo.HostSystem
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		Log.Errorln("[esxi.go]", err)
	}
	defer v.Destroy(ctx)

	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary", "datastore"}, &esxiList)
	if err != nil {
		Log.Errorln("[esxi.go]", err)
	}
	return esxiList
}

//esxiAlive 电源状态
func esxiPower(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue {
	/*
		poweredOff:关机
		poweredOn:开机
		standBy:待机
		unknown:主机断开连接或者无响应时被标记为未知
	*/
	switch esxi.Summary.Runtime.PowerState {
	case "poweredOff":
		return []*MetricValue{GaugeValue("agent.power", "1.0")}
	case "poweredOn":
		return []*MetricValue{GaugeValue("agent.power", "2.0")}
	case "standBy":
		return []*MetricValue{GaugeValue("agent.power", "3.0")}
	case "unknown":
		return []*MetricValue{GaugeValue("agent.power", "4.0")}
	default:
		return []*MetricValue{GaugeValue("agent.power", "4.0")}
	}
}

//esxiStatus 主机状态
func esxiStatus(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue {
	/*
		gray:状态未知
		green:实体没问题
		red:实体肯定有问题
		yellow:实体可能有问题
	*/
	switch esxi.Summary.OverallStatus {
	case "gray":
		return []*MetricValue{GaugeValue("agent.status", "1.0")}
	case "green":
		return []*MetricValue{GaugeValue("agent.status", "2.0")}
	case "red":
		return []*MetricValue{GaugeValue("agent.status", "3.0")}
	case "yellow":
		return []*MetricValue{GaugeValue("agent.status", "4.0")}
	default:
		return []*MetricValue{GaugeValue("agent.status", "1.0")}
	}
}

//esxiUptime 开机时间
func esxiUptime(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue {
	return []*MetricValue{GaugeValue("agent.uptime", esxi.Summary.QuickStats.Uptime)}
}

//esxiCPU CPU相关监控信息
func esxiCPU(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue {
	var total = int64(esxi.Summary.Hardware.CpuMhz) * int64(esxi.Summary.Hardware.NumCpuCores)
	totalCPU := GaugeValue("cpu.total", total*1000*1000)
	useCPU := GaugeValue("cpu.usage.average", int64(esxi.Summary.QuickStats.OverallCpuUsage)*1000*1000)
	usePercentCPU := GaugeValue("cpu.busy", fmt.Sprintf("%.2f", float64(esxi.Summary.QuickStats.OverallCpuUsage)/float64(total)*100))
	freeCPU := GaugeValue("cpu.free.average", (int64(total)-int64(esxi.Summary.QuickStats.OverallCpuUsage))*1000*1000)
	return []*MetricValue{totalCPU, freeCPU, useCPU, usePercentCPU}
}

//esxiMem 内存相关监控信息
func esxiMem(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue {
	var total = esxi.Summary.Hardware.MemorySize
	var free = int64(esxi.Summary.Hardware.MemorySize) - (int64(esxi.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
	totalMem := GaugeValue("mem.memtotal", total)
	useMem := GaugeValue("mem.memused", int64(esxi.Summary.QuickStats.OverallMemoryUsage)*1024*1024)
	freeMem := GaugeValue("mem.memfree", free)
	freeMemPer := GaugeValue("mem.memfree.percent", fmt.Sprintf("%.2f", float64(free)/float64(total)*100))
	usedMemPer := GaugeValue("mem.memused.percent", fmt.Sprintf("%.2f", float64(esxi.Summary.QuickStats.OverallMemoryUsage)/float64(total)*100))
	return []*MetricValue{totalMem, useMem, freeMem, freeMemPer, usedMemPer}
}

//esxiNet 网络相关监控信息
func esxiNet(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue {
	var netPerf []*MetricValue
	counterNameID := CoID()
	var EsxiNetExtend = []string{"net.bytesRx.average", "net.bytesTx.average"}
	extendID := CounterIDByName(counterNameID, EsxiNetExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags string
				if each.Instance != "" {
					tags = "dev=" + each.Instance
				}
				netPerf = append(netPerf, GaugeValue(each.Metric, each.Value, tags))
			}
		}
	}
	return netPerf
}

//esxiDatastore 存储相关监控信息
func esxiDatastore(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue {
	var datastorePerf []*MetricValue
	counterNameID := CoID()
	var EsxiDatastoreExtend = []string{"datastore.totalWriteLatency.average", "datastore.totalReadLatency.average"}
	extendID := CounterIDByName(counterNameID, EsxiDatastoreExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags string
				if each.Instance != "" {
					for _, eachDs := range *dsWURL {
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
	return datastorePerf
}

//esxiDisk 磁盘使用空间
func esxiDisk(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue {
	var diskPerf []*MetricValue
	pc := property.DefaultCollector(c.Client)
	dss := []mo.Datastore{}
	err := pc.Retrieve(ctx, esxi.Datastore, []string{"summary"}, &dss)
	if err != nil {
		Log.Warnln("[esxi.go] get host datastore info err:", err)
	}
	var (
		freeAll  int64
		totalAll int64
		usedAll  int64
	)
	for _, ds := range dss {
		var tags = "datastore=" + ds.Summary.Name
		var free = ds.Summary.FreeSpace
		var total = ds.Summary.Capacity
		var used = total - free
		freeAll += free
		totalAll += total
		usedAll += used
		var freePercent = float64(free) / float64(total) * 100
		var usedPercent = float64(used) / float64(total) * 100
		diskPerf = append(diskPerf, GaugeValue("df.bytes.free", free, tags))
		diskPerf = append(diskPerf, GaugeValue("df.bytes.total", total, tags))
		diskPerf = append(diskPerf, GaugeValue("df.bytes.used", used, tags))
		diskPerf = append(diskPerf, GaugeValue("df.bytes.free.percent", freePercent, tags))
		diskPerf = append(diskPerf, GaugeValue("df.bytes.used.Percent", usedPercent, tags))
	}
	freeAllPercent := float64(freeAll) / float64(totalAll) * 100
	usedAllPercent := float64(usedAll) / float64(totalAll) * 100
	diskPerf = append(diskPerf, GaugeValue("df.statistics.total", totalAll))
	diskPerf = append(diskPerf, GaugeValue("df.statistics.used", usedAll))
	diskPerf = append(diskPerf, GaugeValue("df.statistics.used.percent", usedAllPercent))
	diskPerf = append(diskPerf, GaugeValue("df.statistics.free", freeAll))
	diskPerf = append(diskPerf, GaugeValue("df.statistics.free.percent", freeAllPercent))
	return diskPerf
}

//EsxiMappers Esxi的Map对象
func EsxiMappers() []EFuncsAndInterval {
	interval := Config().Transfer.Interval
	mappers := []EFuncsAndInterval{
		{
			Fs: []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue{
				esxiPower,
				esxiStatus,
				esxiUptime,
				esxiCPU,
				esxiMem,
				esxiNet,
				esxiDatastore,
				esxiDisk,
			},
			Interval: interval,
		},
	}
	return mappers
}
