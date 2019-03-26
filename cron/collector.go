package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vmware/govmomi/vim25/mo"

	"github.com/Vsphere/g"
	"github.com/vmware/govmomi"
)

//Collect 监控项采集
func Collect(ctx context.Context, c *govmomi.Client, cfg *g.VsphereConfig) {
	stime := time.Now().Unix()

	if !g.Config().Transfer.Enabled || len(g.Config().Transfer.Addr) == 0 {
		return
	}

	var wg sync.WaitGroup
	mapperEsxi := g.EsxiMappers()
	mapperVS := g.VsphereMappers()
	esxiList := g.EsxiList(ctx, c)

	g.ReportVCStatus(cfg)
	g.CounterWithID(ctx, c)
	g.DsWithURL(ctx, c)

	for _, v := range mapperVS {
		wg.Add(1)
		go collectVsphere(ctx, c, v, cfg, &wg)
	}

	if cfg.Extend {
		mapperExtend := g.EsxiExtendMappers()
		mapperEsxi = append(mapperEsxi, mapperExtend...)
		for _, esxi := range esxiList {
			g.ReportVCStatus(cfg, esxi.Summary.Config.Name)
		}
	}

	for _, esxi := range esxiList {
		for _, ve := range mapperEsxi {
			wg.Add(1)
			go collectEsxi(ctx, c, esxi, ve, cfg, &wg)
		}
	}

	wg.Wait()
	etime := time.Now().Unix()
	g.Log.Infof("[collector.go] the vc %s have been collected, time:%d", g.InitVCIP(cfg), etime-stime)
}

func collectVsphere(ctx context.Context, c *govmomi.Client, v g.VFuncsAndInterval, cfg *g.VsphereConfig, wg *sync.WaitGroup) {

	defer wg.Done()
	hostname := cfg.Hostname
	mvs := []*g.MetricValue{}

	for _, fn := range v.Fs {
		items := fn(ctx, c)
		if items == nil || len(items) == 0 {
			continue
		}
		for _, mv := range items {
			mvs = append(mvs, mv)
		}
	}

	now := time.Now().Unix()
	for j := 0; j < len(mvs); j++ {
		mvs[j].Step = int64(v.Interval)
		mvs[j].Timestamp = now
		mvs[j].Endpoint = hostname
	}
	g.SendToTransfer(mvs)
}

func collectEsxi(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, v g.EFuncsAndInterval, cfg *g.VsphereConfig, wg *sync.WaitGroup) {

	defer wg.Done()
	mvs := []*g.MetricValue{}
	now := time.Now().Unix()

	mvs = append(mvs, g.AgentMetrics(ctx, c)...)

	for _, fn := range v.Fs {
		items := fn(ctx, c, esxi)
		if items == nil || len(items) == 0 {
			continue
		}
		for _, mv := range items {
			mvs = append(mvs, mv)
		}
	}

	if !cfg.Split {
		for _, x := range mvs {
			tags := fmt.Sprintf("host=%s", esxi.Summary.Config.Name)
			x.Tags = fmt.Sprintf("%s,%s", x.Tags, tags)
			x.Step = int64(v.Interval)
			x.Metric = cfg.MetricHead + x.Metric
			x.Timestamp = now
			x.Endpoint = cfg.Hostname
		}
	} else {
		for _, x := range mvs {
			x.Endpoint = cfg.EndpointHead + esxi.Summary.Config.Name
			x.Step = int64(v.Interval)
			x.Timestamp = now
		}
	}

	g.SendToTransfer(mvs)
}
