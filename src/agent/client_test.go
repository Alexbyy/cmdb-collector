package agent

import (
	"fmt"
	"k8s.io/klog/v2"
	"strings"
	"testing"
)

func TestClient_AddInstance(t *testing.T) {
	fmt.Printf("tehere: \n")
	agent := NewClient("http://10.110.19.61:32033")
	nodes := []Node{
		{Name: "mgt01"},
		{Name: "slave01"},
		{Name: "compute01"},
		{Name: "storage01"},
	}

	for _, node := range nodes {
		var res map[string]interface{}
		var err error
		if strings.HasPrefix(node.Name, "mgt") {
			fmt.Printf("tehere: %v\n", node.Name)
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
		}
		if res["bk_error_msg"] != "success" {
			klog.Infof("addInstance pod %v error, error info: %v\n", res["bk_inst_name"], res["bk_error_msg"])
			continue
		}

	}
}

func TestClient_GetInstanceList(t *testing.T) {
	agent := NewClient("http://10.110.19.61:32033")
	body := make(map[string]interface{})
	res, err := agent.GetInstanceList("/api/v3/find/instassociation/object/pod", body)
	if err != nil {
		klog.Errorf("err %v\n", err)
	}
	for _, cal := range (*res).([]interface{}) {
		fmt.Printf("item %v\n", cal.(map[string]interface{})["bk_inst_id"])
	}
}
