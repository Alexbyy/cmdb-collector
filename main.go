package main

import (
	a "cmdb-collector/src/agent"
	c "cmdb-collector/src/collector"
	m "cmdb-collector/src/manager"
	"fmt"
	"k8s.io/klog/v2"
	"time"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	collector := c.NewCollector()
	agent := a.NewClient("http://10.110.19.61:32033")
	manager := m.NewManager(agent, collector)


	//清理旧实例
	err := agent.ClearData()
	if err != nil{
		klog.Errorf("清理旧数据发生错误： %v\n", err)
	}

	//step1:获取Objects
	objects, err := agent.GetModels()
	if err != nil {
		klog.Errorf("get Objects error: v%\n", err)
	}
	//records := map[string]interface{}{}

	//step2:遍历获取到的objects
	for _, value := range objects["data"].([]interface{}) {
		//此id为模型分组id
		classificationId := value.(map[string]interface{})["bk_classification_id"].(string)
		if classificationId == "bk_host_manage" || classificationId == "bk_biz_topo" || classificationId == "bk_organization" || classificationId == "bk_network" {
			continue
		}


		//获取分组id下的object的数据
		objects := value.(map[string]interface{})["bk_objects"].([]interface{})
		if len(objects) == 0 {
			continue
		}
		for _, item := range objects {
			objId := item.(map[string]interface{})["bk_obj_id"].(string)
			data, err := collector.GetObjData(objId)
			if err != nil || data == nil || len(*data) == 0{
				klog.Errorf("获取object数据:%v 错误：%v, 获取结果： %v\n", objId, err, data)
				continue
			}

			//step3：遍历data，创建实例
			for _, item := range *data {

				switch item := item.(type) {
				case nil:
					continue
				default:
					fmt.Printf("类型为：%T\n", item)
					InstanceRes, err := agent.AddInstance("POST", objId, item)
					if err != nil || InstanceRes["bk_error_msg"] != "success"{
						if InstanceRes["bk_error_code"] == "1199014"{
							// 实例id重复,删除已存在实例后重新创建
							instCon := a.InstCondition{
								Field:    "id",
								Operator: "$eq",
								Value:    item.(map[string]string)["id"],
							}
							var temp []a.InstCondition
							temp = append(temp, instCon)
							condition := a.Condition{
								Condition: map[string]interface{}{
									"pod": temp,
								},
							}

							// 根据条件获取实例的bk_inst_id
							res, _ := agent.GetInstance(objId, &condition)
							bk_inst_id := (res["data"].(map[string]interface{})["info"]).(map[string]interface{})["bk_inst_id"].(string)


							//删除实例
							_, _ = agent.DelInstance(objId, bk_inst_id)

							//重新创建
							InstanceRes, err = agent.AddInstance("POST", objId, item)

						}
						if err != nil{
							klog.Errorf("AddInstance error: %v,AddInstance result: %v\n", err, InstanceRes)
							continue
						}

					}

					//step4：已成功创建实例，这一步获取obj_item的关系
					associRes1, associRes2, err := agent.GetObjAssociation(objId)
					if err != nil || associRes1["bk_error_msg"] != "success" || associRes2["bk_error_msg"] != "success" {
						klog.Errorf("GetObjAssociation error: %v\n", err)
						klog.Errorf("GetObjAssociation  %v error, error info: %v;%v\n", item, associRes1["bk_error_msg"],associRes2["bk_error_msg"])
						continue

					}

					//step5:遍历associRes建立实例间关系
					associ1 := associRes1["data"].([]interface{})
					associ2 := associRes2["data"].([]interface{})
					manager.BuildAssociation(&InstanceRes, &associ1, objId)
					manager.BuildAssociation(&InstanceRes, &associ2, objId)
					fmt.Printf("associ1 len: %v;associ2 len :%v\n", len(associ1), len(associ2))
				}
			}
		}
	}

	for {
		fmt.Println("无限循环中")
		time.Sleep(15 * time.Second)
	}

}
