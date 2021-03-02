package main

import (
	a "cmdb-collector/src/agent"

	c "cmdb-collector/src/collector"
	//"time"
)

func main() {
	collector := c.NewCollector()
	agent := a.NewClient("http://10.110.19.61:32033")

	podList := collector.GetPods("monitoring")
	for _, item := range *podList{
		agent.AddInstance("POST", "pod", item)
	}
	//pod := a.Pods{Name: "test"}
	//agent.AddInstance("POST", "pod", pod)
	//
	//for {
	//
	//
	//
	//	time.Sleep(10 * time.Second)
	//}
}