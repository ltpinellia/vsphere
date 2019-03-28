package g

import (
	"context"
	"net/rpc"
	"sync"
	"time"

	"github.com/vmware/govmomi/vim25/mo"

	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
)

const (
	//VERSION 版本号
	VERSION = "1.0.0"
	//LOGFILE 日志文件
	LOGFILE = "./var/vsphere.log"
)

//Log 全局log
var Log = logrus.New()

//Root 程序根目录
var Root string

//OldModTime 配置文件当前修改时间
var OldModTime int64

//OldExtendTime 旧扩展配置文件修改时间
var OldExtendTime int64

//LocalIP 本地IP
var LocalIP string

//HbsClient HBS客户端
var HbsClient *SingleConnRPCClient

//SimpleRPCResponse RPC返回
type SimpleRPCResponse struct {
	Code int `json:"code"`
}

//SingleConnRPCClient RPC client
type SingleConnRPCClient struct {
	sync.Mutex
	rpcClient *rpc.Client
	RPCServer string
	Timeout   time.Duration
}

//AgentReportRequest Agent上报状态信息格式
type AgentReportRequest struct {
	Hostname      string
	IP            string
	AgentVersion  string
	PluginVersion string
}

//MetricValue 上传监控项格式
type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

//VFuncsAndInterval Vsphere上传监控项函数及上传间隔
type VFuncsAndInterval struct {
	Fs       []func(ctx context.Context, c *govmomi.Client) []*MetricValue
	Interval int
}

//EFuncsAndInterval Esxi上传监控项函数及上传间隔
type EFuncsAndInterval struct {
	Fs       []func(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, dsWURL *[]DatastoreWithURL) []*MetricValue
	Interval int
}

//DatastoreWithURL 存储器与路径之间的映射
type DatastoreWithURL struct {
	Datastore string
	URL       string
}

//TransferResponse transfer返回格式
type TransferResponse struct {
	Message string
	Total   int
	Invalid int
	Latency int64
}

//MetricPerf Performance监控项值
type MetricPerf struct {
	Metric   string `json:"metric"`
	Value    int64  `json:"value"`
	Instance string `json:"instance"`
}
