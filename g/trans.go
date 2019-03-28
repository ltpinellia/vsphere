package g

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func (trs *TransferResponse) String() string {
	return fmt.Sprintf(
		"<Total=%v, Invalid:%v, Latency=%vms, Message:%s>",
		trs.Total,
		trs.Invalid,
		trs.Latency,
		trs.Message,
	)
}

var (
	//TransferClientsLock TransferClientsLock
	TransferClientsLock = new(sync.RWMutex)
	//TransferClients TransferClients
	TransferClients = map[string]*SingleConnRPCClient{}
)

//SendMetrics send metrics
func SendMetrics(metrics []*MetricValue, resp *TransferResponse) {
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(Config().Transfer.Addr)) {
		addr := Config().Transfer.Addr[i]
		c := getTransferClient(addr)
		if c == nil {
			c = initTransferClient(addr)
		}

		if updateMetrics(c, metrics, resp) {
			break
		}
	}
}

func initTransferClient(addr string) *SingleConnRPCClient {
	var c = &SingleConnRPCClient{
		RPCServer: addr,
		Timeout:   time.Duration(Config().Transfer.Timeout) * time.Millisecond,
	}
	TransferClientsLock.Lock()
	defer TransferClientsLock.Unlock()
	TransferClients[addr] = c

	return c
}

func updateMetrics(c *SingleConnRPCClient, metrics []*MetricValue, resp *TransferResponse) bool {
	err := c.Call("Transfer.Update", metrics, resp)
	if err != nil {
		Log.Warnln("[trans.go] call Transfer.Update fail:", c, err)
		return false
	}
	return true
}

func getTransferClient(addr string) *SingleConnRPCClient {
	TransferClientsLock.RLock()
	defer TransferClientsLock.RUnlock()

	if c, ok := TransferClients[addr]; ok {
		return c
	}
	return nil
}

//SendToTransfer send metrics to transfer
func SendToTransfer(metrics []*MetricValue) {
	if len(metrics) == 0 {
		return
	}
	Log.Debugf("[trans.go] => <Total=%d> %v", len(metrics), metrics[0])

	var resp TransferResponse
	SendMetrics(metrics, &resp)
	Log.Debugln("[trans.go] <=", &resp)
}
