package g

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
)

//InitVCIP get vc ip
func InitVCIP(vc *VsphereConfig) string {
	var ip string
	if vc.IP != "" {
		ip = vc.IP
	} else {
		re := regexp.MustCompile(`([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})`)
		match := re.FindAllStringSubmatch(vc.Addr, -1)
		if match != nil {
			ip = match[0][1]
		}
	}
	return ip
}

//InitRPCClients initialize rpc client
func InitRPCClients() {
	if Config().Heartbeat.Enabled {
		HbsClient = &SingleConnRPCClient{
			RPCServer: Config().Heartbeat.Addr,
			Timeout:   time.Duration(Config().Heartbeat.Timeout) * time.Millisecond,
		}
	}
}

//NewMetricValue decorate metric object,return new metric with tags
func NewMetricValue(metric string, val interface{}, dataType string, tags ...string) *MetricValue {
	mv := MetricValue{
		Metric: metric,
		Value:  val,
		Type:   dataType,
	}

	size := len(tags)

	if size > 0 {
		mv.Tags = strings.Join(tags, ",")
	}

	return &mv
}

//GaugeValue Gauge type
func GaugeValue(metric string, val interface{}, tags ...string) *MetricValue {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

//CounterValue Counter type
func CounterValue(metric string, val interface{}, tags ...string) *MetricValue {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}

var (
	coID   *map[string]int32
	dsWUrl *[]DatastoreWithURL
)

//CoID return all esxi Performance。example:disk.io.write_requests:524
func CoID() map[string]int32 {
	lock.RLock()
	defer lock.RUnlock()
	return *coID
}

//DsWURL return dsWUrl
func DsWURL() *[]DatastoreWithURL {
	lock.RLock()
	defer lock.RUnlock()
	return dsWUrl
}

//CounterWithID get all counter list(Key:counter name,Value:counter id)
func CounterWithID(ctx context.Context, c *govmomi.Client) {
	CounterNameID := make(map[string]int32)
	m := performance.NewManager(c.Client)
	p, err := m.CounterInfoByKey(ctx)
	if err != nil {
		Log.Infoln("[common.go]", err)
	}
	for _, cc := range p {
		CounterNameID[cc.Name()] = cc.Key
	}
	coID = &CounterNameID
}

//DsWithURL get all datastore with URL
func DsWithURL(ctx context.Context, c *govmomi.Client) {
	dss := datastores(ctx, c)
	datastoreWithURL := []DatastoreWithURL{}
	if dss != nil {
		for _, ds := range dss {
			datastoreWithURL = append(datastoreWithURL, DatastoreWithURL{Datastore: ds.Summary.Name, URL: ds.Summary.Url})
		}
	}
	dsWUrl = &datastoreWithURL
}

//CounterIDByName get counter key by counter name
func CounterIDByName(CounterNameID map[string]int32, Name []string) []int32 {
	IDList := make([]int32, 0)
	for _, eachName := range Name {
		ID, exit := CounterNameID[eachName]
		if exit {
			IDList = append(IDList, ID)
		}
	}
	return IDList
}
