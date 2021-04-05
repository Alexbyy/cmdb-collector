package main

import (
	a "cmdb-collector/src/agent"
	m "cmdb-collector/src/manager"
	"cmdb-collector/src/options"
	"github.com/robfig/cron"
	"flag"
	"k8s.io/klog/v2"
	"time"
)



func init() {
	// Default logging verbosity to V(2)
	flag.Set("v", "6")
}

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()
	opts := options.NewOptions()

	agent := a.NewClient("http://10.110.19.61:32033")
	manager, err := m.NewManager(agent, opts)
	klog.Infof("Options: %v\n", opts)
	if err != nil {
		klog.Fatalf("初始化NewManager报错:%v\n", err)
	}
	manager.Start()

	//定时任务
	i := 0
	c := cron.New()
	c.AddFunc(opts.Cron, func() {
		i++
		klog.Infof("cron running:%s; time:%s\n", i, time.Now())
		manager.Start()
	})
	c.Start()
	select{}

}



