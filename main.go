package main

import (
	a "cmdb-collector/src/agent"
	c "cmdb-collector/src/collector"
	"fmt"
	"k8s.io/klog/v2"
	"strings"
	"time"
	//"time"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	collector := c.NewCollector()
	agent := a.NewClient("http://10.110.19.61:32033")
	////
	recordPod := map[string]int{}

	//获取数据
	podList, containerList := collector.GetPods("monitoring")
	nodeList, err := collector.GetNodes()
	if err != nil {
		klog.Errorf("GetNodes error: %v\n", err)
	}

	//上传数据
	for _, item := range *podList {
		res, err := agent.AddInstance("POST", "pod", item)
		if err != nil {
			klog.Error("AddInstance pod: %v, error: %v\n", item.Name, err)

		}
		if res["bk_error_msg"] != "success" {
			klog.Infof("addInstance pod %v error, error info: %v\n", item.Name, res["bk_error_msg"])
			continue
		}
		data := res["data"].(map[string]interface{})
		recordPod[item.Name] = int(data["bk_inst_id"].(float64))
	}

	for podname, containers := range *containerList {
		for _, con := range containers {
			res, err := agent.AddInstance("POST", "container", con)
			if err != nil {
				klog.Error("AddInstance container:%v, error: %v\n", con.Name, err)
			}
			if res["bk_error_msg"] != "success" {
				klog.Infof("addInstance container %v error, error info: %v\n", con.Name, res["bk_error_msg"])
				continue
			}
			data := res["data"].(map[string]interface{})
			var1 := int(data["bk_inst_id"].(float64))
			var2 := recordPod[podname]
			body := map[string]interface{}{"bk_asst_inst_id": var1, "bk_inst_id": var2, "bk_obj_asst_id": "pod_group_container"}
			res, err = agent.InstAssoci("POST", "/api/v3/create/instassociation", body)
			if err != nil {
				klog.Error("InstAssoci container error: %v\n",  err)
			}
			if res["bk_error_msg"] != "success" {
				klog.Infof("InstAssoci container error, error info: %v\n", res["bk_error_msg"])
				continue
			}
		}
	}

	recordNode := make(map[string]int)
	klog.Infof("nodelist 的长度为： %s\n", len(nodeList))
	for _, node := range nodeList {
		var res map[string]interface{}
		var err error

		if node == nil || node.Name == "" {
			klog.Infof("node 数据异常：%v", node)
			continue
		}

		if strings.HasPrefix(node.Name, "mgt") {
			res, err = agent.AddInstance("POST", "mgt", node)
		}
		if strings.HasPrefix(node.Name, "slave") {
			res, err = agent.AddInstance("POST", "slave", node)
		}
		if strings.HasPrefix(node.Name, "compute") {
			res, err = agent.AddInstance("POST", "compute", node)
		}
		if strings.HasPrefix(node.Name, "storage") {
			res, err = agent.AddInstance("POST", "storage", node)
		}

		if err != nil {
			klog.Errorf("add instance:%v  node error: %v\n", node.Name, err)
			continue
		}
		if res["bk_error_msg"] != "success" {
			klog.Infof("addInstance node %v error, error info: %v\n", node.Name, res["bk_error_msg"])
			continue
		} else {
			data := res["data"].(map[string]interface{})
			recordNode[node.Name] = int(data["bk_inst_id"].(float64))
		}
	}
	body := make(map[string]interface{})
	podInstanceList, err := agent.GetInstanceList("/api/v3/find/instassociation/object/pod", body)
	if err != nil {
		klog.Errorf("获取podlist err: %v\n", err)
	}

	klog.Infof("recordNode %v\n", recordNode)
	klog.Infof("recordPod: %v\n", podInstanceList)

	//建立pod与node的关联关系
	for _, pod := range (*podInstanceList).([]interface{}) {
		nodename := pod.(map[string]interface{})["icp_pod_nodename"].(string)
		str := ""
		if strings.HasPrefix(nodename, "mgt") {
			str = "mgt"
		}
		if strings.HasPrefix(nodename, "slave") {
			str = "slave"
		}
		if strings.HasPrefix(nodename, "compute") {
			str = "compute"
		}
		if strings.HasPrefix(nodename, "storage") {
			str = "storage"
		}
		if str == "" {
			klog.Infof("str is 空字符串")
			continue
		}

		var nodeId int
		if _, ok := recordNode[nodename]; ok {
			//存在
			nodeId = recordNode[nodename]
		} else {
			continue
		}
		podId := int(pod.(map[string]interface{})["bk_inst_id"].(float64))
		klog.Infof("pod.NodeName:%v, nodeId:%v\n", nodename, nodeId)
		klog.Infof("pod.NodeName:%v, podId%v\n", nodename, pod.(map[string]interface{})["bk_inst_id"])
		body := map[string]interface{}{"bk_asst_inst_id": podId, "bk_inst_id": nodeId, "bk_obj_asst_id": str + "_run_pod"}
		agent.InstAssoci("POST", "/api/v3/create/instassociation", body)
	}

	for {
		fmt.Println("无限循环中")
		time.Sleep(15 * time.Second)
	}

}
