package main

import (
	a "cmdb-collector/src/agent"
	c "cmdb-collector/src/collector"
	"fmt"
	"k8s.io/klog/v2"
	"time"
	//"time"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	collector := c.NewCollector()
	agent := a.NewClient("http://10.110.19.61:32033")

	//step1:获取Objects
	objects, err := agent.GetModels()
	if err != nil {
		klog.Errorf("get Objects error: v%\n", err)
	}
	records := map[string]interface{}{}

	//step2:遍历获取到的objects
	for _, value := range objects["data"].([]map[string]interface{}) {
		//此id为模型分组id
		classificationId := value["bk_classification_id"].(string)
		if classificationId == "bk_host_manage" || classificationId == "bk_biz_topo" || classificationId == "icon-cc-organization" || classificationId == "icon-cc-network-equipment" {
			continue
		}

		//记录分组信息
		if _, ok := records[classificationId]; !ok {
			records[classificationId] = []string{}
		}

		//获取分组id下的object的数据
		objects := value["bk_objects"].([]map[string]interface{})
		for _, item := range objects {
			objId := item["bk_obj_id"].(string)
			data, err := collector.GetObjData(objId)
			if err != nil {
				klog.Errorf("获取object数据:%v 错误：%v\n", objId, err)
			}

			//遍历data，创建实例
			for _, item := range *data {
				//比如pod返回的数据
				if data, ok := item.(struct{}); ok {
					agent.AddInstance("POST", objId, data)
				}

				//比如container返回的数据
				if data, ok := item.([]struct{}); ok {
					for _, data := range data {
						res, err := agent.AddInstance("POST", objId, data)
						if err != nil {
							klog.Errorf("AddInstance error: %v\n", err)

						}
						if res["bk_error_msg"] != "success" {
							klog.Errorf("addInstance  %v error, error info: %v\n", data, res["bk_error_msg"])
							continue
						}
					}
				}
			}

		}

	}

	//////
	//recordPod := map[string]int{}
	//recordSts := map[string]int{}
	//
	////获取数据
	//podList, containerList := collector.GetPods("monitoring")
	//nodeList, err := collector.GetNodes()
	//if err != nil {
	//	klog.Errorf("GetNodes error: %v\n", err)
	//}
	//
	////上传数据
	//for _, item := range *podList {
	//	res, err := agent.AddInstance("POST", "pod", item)
	//	if err != nil {
	//		klog.Errorf("AddInstance pod: %v, error: %v\n", item.Name, err)
	//
	//	}
	//	if res["bk_error_msg"] != "success" {
	//		klog.Errorf("addInstance pod %v error, error info: %v\n", item.Name, res["bk_error_msg"])
	//		continue
	//	}
	//	data := res["data"].(map[string]interface{})
	//	recordPod[item.Name] = int(data["bk_inst_id"].(float64))
	//}
	//
	//fmt.Printf("container list length: %s/n", len(*containerList))
	//fmt.Printf("containerlist: %v\n", *containerList)
	//for podname, containers := range *containerList {
	//	fmt.Printf("container list length: %s\n", len(containers))
	//	for _, con := range containers {
	//		res, err := agent.AddInstance("POST", "container", con)
	//		if err != nil {
	//			klog.Errorf("AddInstance container:%v, error: %v\n", con.Name, err)
	//		}
	//		if res["bk_error_msg"] != "success" {
	//			klog.Errorf("addInstance container %v error, error info: %v, body 数据：%v\n", con.Name, res["bk_error_msg"], con)
	//			continue
	//		}
	//		data := res["data"].(map[string]interface{})
	//		var1 := int(data["bk_inst_id"].(float64))
	//		var2 := recordPod[podname]
	//		body := map[string]interface{}{"bk_asst_inst_id": var1, "bk_inst_id": var2, "bk_obj_asst_id": "pod_group_container"}
	//		res, err = agent.InstAssoci("POST", "/api/v3/create/instassociation", body)
	//		if err != nil {
	//			klog.Errorf("InstAssoci container error: %v\n",  err)
	//		}
	//		if res["bk_error_msg"] != "success" {
	//			klog.Errorf("InstAssoci container error, error info: %v\n", res["bk_error_msg"])
	//			continue
	//		}
	//	}
	//}
	//
	//recordNode := make(map[string]int)
	//klog.Infof("nodelist 的长度为： %s\n", len(nodeList))
	//for _, node := range nodeList {
	//	var res map[string]interface{}
	//	var err error
	//
	//	if node == nil || node.Name == "" {
	//		klog.Infof("node 数据异常：%v", node)
	//		continue
	//	}
	//
	//	if strings.HasPrefix(node.Name, "mgt") {
	//		res, err = agent.AddInstance("POST", "mgt", node)
	//	}
	//	if strings.HasPrefix(node.Name, "slave") {
	//		res, err = agent.AddInstance("POST", "slave", node)
	//	}
	//	if strings.HasPrefix(node.Name, "compute") {
	//		res, err = agent.AddInstance("POST", "compute", node)
	//	}
	//	if strings.HasPrefix(node.Name, "storage") {
	//		res, err = agent.AddInstance("POST", "storage", node)
	//	}
	//
	//	if err != nil {
	//		klog.Errorf("add instance:%v  node error: %v\n", node.Name, err)
	//		continue
	//	}
	//	if res["bk_error_msg"] != "success" {
	//		klog.Error("addInstance node %v error, error info: %v\n", node.Name, res["bk_error_msg"])
	//		continue
	//	} else {
	//		data := res["data"].(map[string]interface{})
	//		recordNode[node.Name] = int(data["bk_inst_id"].(float64))
	//	}
	//}
	//body := make(map[string]interface{})
	//podInstanceList, err := agent.GetInstanceList("/api/v3/find/instassociation/object/pod", body)
	//if err != nil {
	//	klog.Errorf("获取podlist err: %v\n", err)
	//}
	//
	//klog.Infof("recordNode %v\n", recordNode)
	//klog.Infof("recordPod: %v\n", podInstanceList)
	//
	////建立pod与node的关联关系
	//for _, pod := range (*podInstanceList).([]interface{}) {
	//	nodename := pod.(map[string]interface{})["icp_pod_nodename"].(string)
	//	str := ""
	//	if strings.HasPrefix(nodename, "mgt") {
	//		str = "mgt"
	//	}
	//	if strings.HasPrefix(nodename, "slave") {
	//		str = "slave"
	//	}
	//	if strings.HasPrefix(nodename, "compute") {
	//		str = "compute"
	//	}
	//	if strings.HasPrefix(nodename, "storage") {
	//		str = "storage"
	//	}
	//	if str == "" {
	//		klog.Infof("str is 空字符串")
	//		continue
	//	}
	//
	//	var nodeId int
	//	if _, ok := recordNode[nodename]; ok {
	//		//存在
	//		nodeId = recordNode[nodename]
	//	} else {
	//		continue
	//	}
	//	podId := int(pod.(map[string]interface{})["bk_inst_id"].(float64))
	//	klog.V(3).Infof("pod.NodeName:%v, nodeId:%v\n", nodename, nodeId)
	//	klog.V(3).Infof("pod.NodeName:%v, podId%v\n", nodename, pod.(map[string]interface{})["bk_inst_id"])
	//	body := map[string]interface{}{"bk_asst_inst_id": podId, "bk_inst_id": nodeId, "bk_obj_asst_id": str + "_run_pod"}
	//	agent.InstAssoci("POST", "/api/v3/create/instassociation", body)
	//}
	//
	////上传sts、ds、deploy
	//stsList, err := collector.GetStatefulsets("monitoring")
	//deployList, err := collector.GetDeployments("monitoring")
	//dsList, err := collector.GetDaemonSets("monitoring")
	//if err != nil{
	//	klog.Errorf("获取workload err: %v\n", err)
	//}
	//
	//for _, deploy := range deployList {
	//	res, err := agent.AddInstance("POST", "deploy", deploy)
	//	if err != nil{
	//		klog.Errorf("AddInstance deploy err: %v\n", err)
	//	}
	//	if res["bk_error_msg"] != "success" {
	//		klog.Errorf("addInstance deploy %v error, error info: %v\n", deploy.Name, res["bk_error_msg"])
	//		continue
	//	}
	//}
	//
	//for _, ds := range dsList {
	//	res, err := agent.AddInstance("POST", "ds", ds)
	//	if err != nil{
	//		klog.Errorf("AddInstance deploy err: %v\n", err)
	//	}
	//	if res["bk_error_msg"] != "success" {
	//		klog.Errorf("addInstance deploy %v error, error info: %v\n", ds.Name, res["bk_error_msg"])
	//		continue
	//	}
	//}
	//
	//for _, sts := range stsList {
	//	res, err := agent.AddInstance("POST", "sts", sts)
	//	if err != nil{
	//		klog.Errorf("AddInstance statefulset err: %v\n", err)
	//	}
	//	if res["bk_error_msg"] != "success" {
	//		klog.Errorf("addInstance statefulset %v error, error info: %v\n", sts.Name, res["bk_error_msg"])
	//		continue
	//	} else {
	//		data := res["data"].(map[string]interface{})
	//		recordSts[sts.Name] = int(data["bk_inst_id"].(float64))
	//	}
	//}
	//
	////建立sts实例与pod实例之间的关联
	//for _, pod := range (*podInstanceList).([]interface{}) {
	//	kind := pod.(map[string]interface{})["icp_pod_kind"].(string)
	//	name := pod.(map[string]interface{})["bk_inst_name"].(string)
	//	podId := int(pod.(map[string]interface{})["bk_inst_id"].(float64))
	//	klog.Infof("kind:%v, name:%v, podID:%v\n", kind, name, podId)
	//	if kind == "StatefulSet" {
	//		for key, stsId := range recordSts {
	//			klog.Infof("sts name:%v\n", key)
	//			if strings.HasPrefix(name, key){
	//				body := map[string]interface{}{"bk_asst_inst_id": podId, "bk_inst_id": stsId, "bk_obj_asst_id": "sts_create_pod"}
	//				res, err := agent.InstAssoci("POST", "/api/v3/create/instassociation", body)
	//				if err != nil{
	//					klog.Errorf("InstAssoci statefulset err: %v\n", err)
	//				}
	//				if res["bk_error_msg"] != "success" {
	//					klog.Errorf("InstAssoci error info: %v\n", res["bk_error_msg"])
	//					continue
	//				}
	//			}
	//		}
	//	}
	//}
	//

	for {
		fmt.Println("无限循环中")
		time.Sleep(15 * time.Second)
	}

}
