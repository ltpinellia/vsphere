package g

import (
	"context"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
)

//InitRootDir 初始化主目录
func InitRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		Log.Fatalln("[common.go] getwd fail: ", err)
	}
}

//InitVCIP 初始化Vsphere IP
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

//InitRPCClients 初始化rpc客户端
func InitRPCClients() {
	if Config().Heartbeat.Enabled {
		HbsClient = &SingleConnRPCClient{
			RPCServer: Config().Heartbeat.Addr,
			Timeout:   time.Duration(Config().Heartbeat.Timeout) * time.Millisecond,
		}
	}
}

//NewMetricValue 返回新的metric对象，主要修改了tags跟type
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

//GaugeValue Gauge型监控项
func GaugeValue(metric string, val interface{}, tags ...string) *MetricValue {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

//CounterValue Counter型监控项
func CounterValue(metric string, val interface{}, tags ...string) *MetricValue {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}

var (
	coID   *map[string]int32
	dsWUrl *[]DatastoreWithURL
)

//CoID 返回esxi Performance监控项名及其Key之间的映射,类似:disk.io.write_requests:524
func CoID() map[string]int32 {
	lock.RLock()
	defer lock.RUnlock()
	return *coID
}

//DsWURL 存储器ID与存储器命名之间的映射
func DsWURL() *[]DatastoreWithURL {
	lock.RLock()
	defer lock.RUnlock()
	return dsWUrl
}

//CounterWithID 获取所有的counter列表(Key:counter name,Value:counter id)
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

//DsWithURL 存储器命名与存储器ID的映射
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

//CounterIDByName 通过counter列表获取Key列表
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
