package main

import (
	a "cmdb-collector/src/agent"
	m "cmdb-collector/src/manager"
	"cmdb-collector/src/options"
	"fmt"
	"k8s.io/klog/v2"
	"time"
)




func main() {
	opts := options.NewOptions()
	opts.AddFlags()
	err := opts.Parse()
	if err != nil {
		klog.Fatalf("Error: %s", err)
	}

	agent := a.NewClient("http://10.110.19.61:32033")
	manager, err := m.NewManager(agent, opts)
	fmt.Printf("Options: %v\n", opts)
	if err != nil {
		klog.Fatalf("初始化NewManager报错:%v\n", err)
	}
	manager.Start()

	for {
		fmt.Println("无限循环中")
		time.Sleep(60 * time.Second)
	}

}



