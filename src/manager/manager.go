package manager

import (
	"cmdb-collector/src/agent"
	"cmdb-collector/src/collector"
	"fmt"
	"k8s.io/klog/v2"
	"strings"
)

type Manager struct {
	Agent *agent.Client
	Collector *collector.Collector
}

func NewManager(a *agent.Client, c *collector.Collector)*Manager  {
	return &Manager{
		Agent: a,
		Collector: c,

	}
}

//建立管理

func  (m *Manager) BuildAssociation(instanceData *map[string]interface{}, associInstList *[]interface{}, objId string){
	associ := *associInstList
	if len(associ) > 0 {
		for _, item := range associ {
			bkObjAsstId := item.(map[string]interface{})["bk_obj_id"].(string) // 要建立关联的id
			bkAsstObjId := item.(map[string]interface{})["bk_obj_asst_id"].(string)  //示例：ds_create_pod
			fmt.Printf("bkObjAsstId:%v;bkAsstObjId:%v\n", bkObjAsstId, bkAsstObjId)
			if _, ok := agent.Association[bkAsstObjId]; !ok{
				klog.Errorf("没有配置实例关系如何配置：bk_asst_obj_id: %v;bk_obj_id:%v; bk_obj_asst_id:%v\n", objId, bkObjAsstId, bkAsstObjId)
				continue
			}
			instanData := (*instanceData)["data"].(map[string]interface{})
			rule := agent.Association[bkAsstObjId]
			temp1 := instanData[rule[objId]].(string)
			instCon := agent.InstCondition{
				Field:    rule[bkObjAsstId],
				Operator: "$eq",
				Value:    temp1,
			}
			var temp []agent.InstCondition
			temp = append(temp, instCon)
			condition := agent.Condition{
				Condition: map[string]interface{}{
					"pod": temp,
				},
			}
			res, err := m.Agent.GetInstance(bkObjAsstId, &condition)

			if err != nil || res["bk_error_msg"] != "success"{
				klog.Errorf("获取实例:  %v,错误：%V\n",bkAsstObjId, err)
				continue
			}
			data := res["data"].(map[string]interface{})["info"]
			for _, item := range data.([]interface{}) {
				var1 := item.(map[string]interface{})["bk_inst_id"].(float64)
				var2 := instanData["bk_inst_id"].(float64)
				body := map[string]interface{}{}
				if strings.HasPrefix(bkAsstObjId, objId){
					body = map[string]interface{}{"bk_asst_inst_id": var1, "bk_inst_id": var2, "bk_obj_asst_id": bkAsstObjId}
				}
				body = map[string]interface{}{"bk_asst_inst_id": var2, "bk_inst_id": var1, "bk_obj_asst_id": bkAsstObjId}
				fmt.Printf("body: %v\n", body)
				res, err := m.Agent.InstAssoci("POST", "/api/v3/create/instassociation", body)
				if err != nil || res["bk_error_msg"] != "success"{
					klog.Errorf("建立实例关系错误：err:%v; obj_item: %v;assObjId:%v; bk_obj_asst_id:%v\n",err, objId, bkAsstObjId, bkObjAsstId)
					continue
				}

			}
		}
	}
}
