package g

import (
	"errors"
	"math"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

func (rp *SingleConnRPCClient) close() {
	if rp.rpcClient != nil {
		rp.rpcClient.Close()
		rp.rpcClient = nil
	}
}

func (rp *SingleConnRPCClient) serverConn() error {
	if rp.rpcClient != nil {
		return nil
	}

	var err error
	var retry = 1

	for {
		if rp.rpcClient != nil {
			return nil
		}

		rp.rpcClient, err = JSONRPCClient("tcp", rp.RPCServer, rp.Timeout)
		if err != nil {
			Log.Warnf("[rpc.go] dial %s fail: %v", rp.RPCServer, err)
			if retry > 3 {
				return err
			}
			time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
			retry++
			continue
		}
		return err
	}
}

//Call RpcClient.call
func (rp *SingleConnRPCClient) Call(method string, args interface{}, reply interface{}) error {

	rp.Lock()
	defer rp.Unlock()

	err := rp.serverConn()
	if err != nil {
		return err
	}

	timeout := time.Duration(10 * time.Second)
	done := make(chan error, 1)

	go func() {
		err := rp.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		Log.Warnf("[rpc.go] rpc call timeout %v => %v", rp.rpcClient, rp.RPCServer)
		rp.close()
		return errors.New(rp.RPCServer + " rpc call timeout")
	case err := <-done:
		if err != nil {
			rp.close()
			return err
		}
	}

	return nil
}

//JSONRPCClient JSONRPCClient
func JSONRPCClient(network, address string, timeout time.Duration) (*rpc.Client, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewClient(conn), err
}
