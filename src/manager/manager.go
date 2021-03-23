package manager

import (
	"cmdb-collector/src/agent"
	"cmdb-collector/src/collector"
	"encoding/json"
	"io/ioutil"
	"k8s.io/klog/v2"
	"strings"
)

type Manager struct {
	Agent *agent.Client
	Collector *collector.Collector
	Config *Config
}

type Config struct {
	AssociConfig map[string]map[string]string `json:"AssociConfig"`
}
type InitParams struct {
	ConfigPath string
}

func NewManager(a *agent.Client, c *collector.Collector, params *InitParams)(*Manager, error)  {
	res, err := ReadFromJson(params.ConfigPath)
	if err != nil {
		return nil, err
	}
	return &Manager{
		Agent: a,
		Collector: c,
		Config: res,

	}, nil
}

//建立管理

func  (m *Manager) BuildAssociation(instanceData *map[string]interface{}, associInstList *[]interface{}, objId string){
	if *instanceData == nil || *associInstList == nil ||(*instanceData)["data"] == nil{
		klog.V(4).Infof("BuildAssociation received nil params: instanceData: %v;associInstList: %v\n", *instanceData, *associInstList)
		return
	}
	klog.V(4).Infof("此时的objId：%v", objId)
	associ := *associInstList
	if len(associ) > 0 {
		for _, item := range associ {

			bkObjAsstId := ""
			bkAsstObjId := item.(map[string]interface{})["bk_obj_asst_id"].(string)  //示例：ds_create_pod
			if strings.HasPrefix(bkAsstObjId, objId){
				bkObjAsstId = item.(map[string]interface{})["bk_asst_obj_id"].(string) // 要建立关联的id
			}
			if strings.HasSuffix(bkAsstObjId, objId){
				bkObjAsstId = item.(map[string]interface{})["bk_obj_id"].(string) // 要建立关联的id
			}
			if bkAsstObjId == objId{
				continue
			}
			//fmt.Printf("objId:%v;bkObjAsstId:%v;bkAsstObjId:%v\n",objId, bkObjAsstId, bkAsstObjId)
			if _, ok := m.Config.AssociConfig[bkAsstObjId]; !ok{
				klog.V(3).Infof("没有配置实例关系如何维护：bk_asst_obj_id: %v;bk_obj_id:%v; bk_obj_asst_id:%v\n", objId, bkObjAsstId, bkAsstObjId)
				continue
			}
			instanData := (*instanceData)["data"].(map[string]interface{})
			rule := m.Config.AssociConfig[bkAsstObjId]

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
					bkObjAsstId: temp,
				},
			}
			klog.V(4).Infof("获取实例的条件： condition:%v;bkObjAsstId: %v", condition,bkObjAsstId)
			res, err := m.Agent.GetInstance(bkObjAsstId, &condition)
			klog.V(4).Infof("获取的实例内容：%v", res)

			if err != nil || res["bk_error_msg"] != "success" {
				klog.Errorf("获取实例:  %v,错误：%V\n",bkAsstObjId, err)
				continue
			}
			data := res["data"].(map[string]interface{})["info"]
			if len(data.([]interface{})) == 0{
				klog.V(3).Infof("依据条件未查询到实例")
			}
			for _, item := range data.([]interface{}) {
				var1 := item.(map[string]interface{})["bk_inst_id"].(float64)
				var2 := instanData["bk_inst_id"].(float64)
				body := map[string]interface{}{}
				if strings.HasPrefix(bkAsstObjId, objId){
					body = map[string]interface{}{"bk_asst_inst_id": var1, "bk_inst_id": var2, "bk_obj_asst_id": bkAsstObjId}
				}else if strings.HasSuffix(bkAsstObjId, objId) {
					body = map[string]interface{}{"bk_asst_inst_id": var2, "bk_inst_id": var1, "bk_obj_asst_id": bkAsstObjId}
				}
				res, err := m.Agent.InstAssoci("POST", "/api/v3/create/instassociation", body)
				if err != nil || res["bk_error_msg"] != "success"{
					klog.Errorf("建立实例关系错误：err:%v; result: %v; obj_item: %v;assObjId:%v; bk_obj_asst_id:%v; body: %v\n",err,res, objId, bkAsstObjId, bkObjAsstId, body)
					continue
				}

			}
		}
	}
}

func ReadFromJson(src string)(*Config, error){
	data,err:= ioutil.ReadFile(src)

	if err != nil{
		klog.V(3).Infof("读取配置文件出错：%v\n", err)
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(data,&config)
	if err != nil{
		klog.V(3).Infof("解析配置文件出错：%v\n", err)
		return nil, err
	}
	return config, nil
}

