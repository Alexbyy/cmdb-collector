package manager

import (
	"cmdb-collector/src/agent"
	"cmdb-collector/src/collector"
	"cmdb-collector/src/options"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"strings"
)

type Manager struct {
	Agent *agent.Client
	Config *Config //建立关系的配置文件
	Options  *options.Options

}

type Config struct {
	AssociConfig map[string]map[string]string `json:"AssociConfig"`
}
type InitParams struct {
	ConfigPath string
}

func NewManager(a *agent.Client, opts *options.Options)(*Manager, error)  {
	res, err := ReadFromJson(opts.ConfigPath)
	if err != nil {
		return nil, err
	}
	return &Manager{
		Agent: a,
		Config: res,
		Options: opts,

	}, nil
}

func (m *Manager) Start(){
	//获取集群信息
	klog.Infof("Get k8s config》》》》》》》》》》》》》》")
	k8s, err := getK8s(m.Options)
	if err != nil {
		klog.Errorf("获取k8s信息出错:%v\n", err)
		return
	}
	m.Options.K8s = k8s


	//清理旧实例
	klog.Infof("Cleaning old data》》》》》》》》》》》》》》")
	err = m.Agent.ClearAllAssociations()
	err = m.Agent.ClearAllInstance()
	if err != nil{
		klog.Fatalf("清理旧数据发生错误： %v\n", err)
	}else {
		klog.V(2).Infof("Cleaning old data done")
	}

	//step1:获取Objects
	klog.Infof("get objects》》》》》》》》》》》》》》")
	objects, err := m.Agent.GetModels()
	if err != nil {
		klog.Fatalf("get Objects error: v%\n", err)
	}

	klog.Infof("一共有%v个集群", len(*k8s))
	for i := 0; i < len(*k8s); i++ {
		fmt.Printf("启动线程：%v\n", i)
		config := ((*k8s)[i]).(map[string]interface{})
		go m.Run(config, objects)
	}
}

func (m *Manager) Run(config map[string]interface{}, objects  map[string]interface{})  {
	collector, err := collector.NewCollector(m.Options, config)
	if err != nil {
		klog.Errorf("创建collector 报错： %v\n", err)
		return
	}

	//step2:遍历获取到的objects
	klog.Infof("遍历objects》》》》》》》》》》》》")
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
					InstanceRes, err := m.Agent.AddInstance("POST", objId, item)
					if err != nil || InstanceRes["bk_error_msg"] != "success"{
						klog.Errorf("AddInstance error: %v,AddInstance result: %v, objId: %v, data: %v\n", err, InstanceRes, objId, item)
						continue
					}

					//step4：已成功创建实例，这一步获取obj_item的关系
					associRes1, associRes2, err := m.Agent.GetObjAssociation(objId)
					if err != nil || associRes1["bk_error_msg"] != "success" || associRes2["bk_error_msg"] != "success" {
						klog.Errorf("GetObjAssociation error: %v\n", err)
						klog.Errorf("GetObjAssociation  %v error, error info: %v;%v\n", item, associRes1["bk_error_msg"],associRes2["bk_error_msg"])
						continue

					}

					//step5:遍历associRes建立实例间关系
					associ1 := associRes1["data"].([]interface{})
					associ2 := associRes2["data"].([]interface{})
					m.BuildAssociation(&InstanceRes, &associ1, objId)
					m.BuildAssociation(&InstanceRes, &associ2, objId)
					//fmt.Printf("associ1 len: %v;associ2 len :%v\n", len(associ1), len(associ2))
				}
			}
		}
	}

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

			temp1 := ""
			if _, ok := instanData[rule[objId]]; ok {
				//存在
				temp1 = instanData[rule[objId]].(string)
			}else {
				klog.Errorf("实例的维护关系key不存在，objId：%v,rule[objId]:%v,instanData:%v\n", objId, rule[objId], instanceData)
				continue
			}

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

func getK8s(opts *options.Options)(*[]interface{},error)  {
	var res []interface{}
	client := &http.Client{}
	url := opts.LcmBaseUrl +  "/lcm/v1/sites/" + opts.LcmSite + "/branches/" + opts.LcmBranch + "/kuberneteses"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	//req.Header.Set("Content-Type", c.ContentType)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	json.Unmarshal(body, &res)
	return &res, nil
}

