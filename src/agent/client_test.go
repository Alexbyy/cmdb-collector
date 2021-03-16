package agent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"k8s.io/klog/v2"
	"strconv"
	"testing"
)

//func TestClient_AddInstance(t *testing.T) {
//	fmt.Printf("tehere: \n")
//	agent := NewClient("http://10.110.19.61:32033")
//	nodes := []Node{
//		{Name: "mgt01"},
//		{Name: "slave01"},
//		{Name: "compute01"},
//		{Name: "storage01"},
//	}
//
//	for _, node := range nodes {
//		var res map[string]interface{}
//		var err error
//		if strings.HasPrefix(node.Name, "mgt") {
//			fmt.Printf("tehere: %v\n", node.Name)
//			res, err = agent.AddInstance("POST", "mgt", node)
//		}
//		if strings.HasPrefix(node.Name, "slave") {
//			res, err = agent.AddInstance("POST", "slave", node)
//		}
//		if strings.HasPrefix(node.Name, "compute") {
//			res, err = agent.AddInstance("POST", "compute", node)
//		}
//		if strings.HasPrefix(node.Name, "storage") {
//			res, err = agent.AddInstance("POST", "storage", node)
//		}
//
//		if err != nil {
//			klog.Errorf("add instance:%v  node error: %v\n", node.Name, err)
//		}
//		if res["bk_error_msg"] != "success" {
//			klog.Infof("addInstance pod %v error, error info: %v\n", res["bk_inst_name"], res["bk_error_msg"])
//			continue
//		}
//
//	}
//}
//
//func TestClient_GetInstanceList(t *testing.T) {
//	agent := NewClient("http://10.110.19.61:32033")
//	body := make(map[string]interface{})
//	res, err := agent.GetInstanceList("/api/v3/find/instassociation/object/pod", body)
//	if err != nil {
//		klog.Errorf("err %v\n", err)
//	}
//	for _, cal := range (*res).([]interface{}) {
//		fmt.Printf("item %v\n", cal.(map[string]interface{})["bk_inst_id"])
//	}
//}

//获取所有实例关系，并删除该关系
func TestClient_GetAssociList(t *testing.T) {
	agent := NewClient("http://10.110.19.61:32033")
	//step1:获取Objects
	objects, err := agent.GetModels()
	if err != nil {
		klog.Errorf("get Objects error: v%\n", err)
	}
	for _, value := range objects["data"].([]interface{}) {
		//此id为模型分组id
		classificationId := value.(map[string]interface{})["bk_classification_id"].(string)
		if classificationId == "bk_host_manage" || classificationId == "bk_biz_topo" || classificationId == "bk_organization" || classificationId == "bk_network" {
			continue
		}

		//获取分组id下的object的数据
		objects := value.(map[string]interface{})["bk_objects"].([]interface{})
		for _, item := range objects {
			objId := item.(map[string]interface{})["bk_obj_id"].(string)
			associRes1, associRes2, err := agent.GetObjAssociation(objId)
			if err != nil {
				klog.Errorf("GetObjAssociation error: %v\n", err)
				continue

			}
			if associRes1["bk_error_msg"] != "success" || associRes2["bk_error_msg"] != "success" {
				klog.Errorf("GetObjAssociation  %v error, error info: %v;%v\n", item, associRes1["bk_error_msg"],associRes2["bk_error_msg"])
				continue
			}
			//step5:遍历associRes建立实例间关系
			associ1 := associRes1["data"].([]interface{})
			if len(associ1) > 0 {
				for _, item := range associ1 {
					bkAsstObjId := item.(map[string]interface{})["bk_obj_asst_id"].(string)  //示例：ds_create_pod
					res, err := agent.GetAssociList(bkAsstObjId)
					klog.Infof("Get Asslist result:%s\n", res["bk_error_msg"])
					if err != nil {
						klog.Errorf("GetAssociList error: %v\n", err)
					}
					for _, item := range res["data"].([]interface{}) {
						id := strconv.Itoa(int(item.(map[string]interface{})["id"].(float64)))
						klog.Infof("id: %s\n", id)
						res, err := agent.DelAssoci(id)
						klog.Infof("del result:%s\n", res["bk_error_msg"])
						if err != nil {
							klog.Errorf("DElAssociList error: %v\n", err)
						}
					}

				}
			}

			associ2 := associRes2["data"].([]interface{})
			if len(associ2) > 0 {
				for _, item := range associ2 {
					bkAsstObjId := item.(map[string]interface{})["bk_obj_asst_id"].(string)  //示例：ds_create_pod
					res, err := agent.GetAssociList(bkAsstObjId)
					klog.Infof("Get Asslist result:%s\n", res["bk_error_msg"])
					if err != nil {
						klog.Errorf("GetAssociList error: %v\n", err)
					}
					for _, item := range res["data"].([]interface{}) {
						id := strconv.Itoa(int(item.(map[string]interface{})["id"].(float64)))
						klog.Infof("id: %s\n", id)
						res, err := agent.DelAssoci(id)
						klog.Infof("del result:%s\n", res["bk_error_msg"])
						if err != nil {
							klog.Errorf("DElAssociList error: %v\n", err)
						}
					}
				}
			}
		}
	}
}

//获取所有实例并删除该实例
func TestClient_DelInstanceList(t *testing.T) {
	agent := NewClient("http://10.110.19.61:32033")


	//step1:获取Objects
	objects, err := agent.GetModels()
	if err != nil {
		klog.Errorf("get Objects error: v%\n", err)
	}
	for _, value := range objects["data"].([]interface{}) {
		//此id为模型分组id
		classificationId := value.(map[string]interface{})["bk_classification_id"].(string)
		if classificationId == "bk_host_manage" || classificationId == "bk_biz_topo" || classificationId == "bk_organization" || classificationId == "bk_network" {
			continue
		}

		//获取分组id下的object的数据
		objects := value.(map[string]interface{})["bk_objects"].([]interface{})
		for _, item := range objects {
			objId := item.(map[string]interface{})["bk_obj_id"].(string)
			res, err := agent.GetInstanceList(objId, nil)
			if err != nil{
				klog.Errorf("获取实例:  %v,错误：%V\n",objId, err)
				continue
			}

			for _, item := range (*res).([]interface{}) {
				bkInstId := int(item.(map[string]interface{})["bk_inst_id"].(float64))
				id := strconv.Itoa(bkInstId)
				res1, err := agent.DelInstance(objId, id)
				klog.Infof("删除实例结果：%v,\n", res1)
				assert.Nil(t , err)
			}
		}

	}

}


//获取所有实例并批量删除该实例
func TestClient_DelInstanceBatch(t *testing.T) {
	agent := NewClient("http://10.110.19.61:32033")

	//step1:获取Objects
	objects, err := agent.GetModels()
	if err != nil {
		klog.Errorf("get Objects error: v%\n", err)
	}
	for _, value := range objects["data"].([]interface{}) {
		//此id为模型分组id
		classificationId := value.(map[string]interface{})["bk_classification_id"].(string)
		if classificationId == "bk_host_manage" || classificationId == "bk_biz_topo" || classificationId == "bk_organization" || classificationId == "bk_network" {
			continue
		}
		var inst_ids []int
		//获取分组id下的object的数据
		objects := value.(map[string]interface{})["bk_objects"].([]interface{})
		for _, item := range objects {
			objId := item.(map[string]interface{})["bk_obj_id"].(string)
			res, err := agent.GetInstanceList(objId, nil)
			if err != nil{
				klog.Errorf("获取实例:  %v,错误：%V\n",objId, err)
				continue
			}

			for _, item := range (*res).([]interface{}) {
				bkInstId := int(item.(map[string]interface{})["bk_inst_id"].(float64))
				inst_ids = append(inst_ids, bkInstId)

			}
			if len(inst_ids) == 0{
				continue
			}

			delInsts := DelInstances{Delete: InstIds{Instids: inst_ids}}
			klog.Infof("delInsts：%v,\n", delInsts)
			res1, err1 := agent.DelInstancesArray(objId, &delInsts)
			klog.Infof("error：%v,\n", err1)
			assert.Nil(t, err1)
			klog.Infof("删除实例结果：%v,\n", res1)

		}
	}
}

func TestClient_GetModels(t *testing.T) {
	agent := NewClient("http://10.110.19.61:32033")
	res, err := agent.GetModels()
	if err != nil {
		klog.Errorf("GetAssociList error: %v\n", err)
	}
	PrintJson(res)

}

func TestClient_GetObjAssociation(t *testing.T) {
	agent := NewClient("http://10.110.19.61:32033")
	ass1, ass2, _ := agent.GetObjAssociation("pod")
	PrintJson(ass1)
	PrintJson(ass2)
	fmt.Printf("ass1 value: %v", ass1["data"])
	fmt.Printf("ass2 value: %v", ass2["data"])

}



func TestClient_InstAssoci(t *testing.T) {
	agent := NewClient("http://10.110.19.61:32033")
	body := map[string]interface{}{"bk_asst_inst_id":4831, "bk_inst_id":4755, "bk_obj_asst_id":"pod_group_container"}

	res, err := agent.InstAssoci("POST", "/api/v3/create/instassociation", body)
	fmt.Printf("err:    %v\n", err)
	fmt.Printf("value:    %v\n", res)

}

//测试根据条件获取某个实例
func TestClient_GetInstance(t *testing.T) {
	fmt.Printf("dwadawdwadawdaw")
	agent := NewClient("http://10.110.19.61:32033")
	objId := "pod"

	instCon := InstCondition{
		Field:    "icp_pod_id",
		Operator: "$eq",
		Value:    "1",
	}
	temp := []InstCondition{}
	temp = append(temp, instCon)
	condition := Condition{
		Condition: map[string]interface{}{
			"pod": temp,
		},
	}
	res, err := agent.GetInstance(objId, &condition)
	fmt.Printf("err:    %v\n", err)
	fmt.Printf("value:    %v\n", res)
}

