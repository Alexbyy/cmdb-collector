package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"strconv"
)

type Client struct {
	BaseUrl     string
	CookieMap   map[string]string
	CookieStr   string
	ContentType string
	HttpClient  *http.Client
	objs        []map[string]string
}

//用于使用条件查询某个实例
type Condition struct{
	Condition map[string]interface{} `json:"condition"`
}
type  InstCondition struct {
	Field string `json:"field"`
	Operator string `json:"operator"`
	Value  string `json:"value"`
}

//用于批量删除实例
type InstIds struct {
	Instids []int `json:"inst_ids"`
}
type DelInstances struct {
	Delete InstIds `json:"delete"`
}

//type DelInstances struct {
//	Delete struct {
//		InstIds []int `json:"inst_ids"`
//	} `json:"delete"`
//}



func NewClient(url string) *Client {
	cookies := map[string]string{
		"HTTP_BLUEKING_SUPPLIER_ID": "0",
		"http_scheme":               "http",
		"cc3":                       "MTYxMTkxMzIzMHxOd3dBTkVGQ1NVTkZTVTVHVTA0MVZUTk5XVkUxVlVSRlVFNUVTRmRMVXpWRFZsSTFTRFpHVGtVeldEZFBRMFEyUjFCRlRVOUpVRkU9fBNp_HYCb_mz_B0U210DQL9ZLcp48P2rA1PPJb3CLwZJ",
	}
	cookieStr := "HTTP_BLUEKING_SUPPLIER_ID=0; http_scheme=http; cc3=MTYxNDU4NDYzNXxOd3dBTkVGQ1NVTkZTVTVHVTA0MVZUTk5XVkUxVlVSRlVFNUVTRmRMVXpWRFZsSTFTRFpHVGtVeldEZFBRMFEyUjFCRlRVOUpVRkU9fOEiES7B0ibwzrTblCC4X45LnXWkD-4_Cax53luOa6C3"
	newClient := &Client{
		BaseUrl:     "http://10.110.19.61:32033",
		ContentType: "application/json;charset=UTF-8",
		CookieMap:   cookies,
		CookieStr:   cookieStr,
		HttpClient:  &http.Client{},
	}
	return newClient
}

func (c *Client) AddInstance(method string, url string, body interface{}) (map[string]interface{}, error) {
	ms, err := json.Marshal(body)
	if err != nil {
		fmt.Errorf("json 编译错误：%v\n", err)
	}
	payload := bytes.NewBuffer([]byte(ms))
	url = c.BaseUrl + "/api/v3/create/instance/object/" + url
	req, err := http.NewRequest(method, url, payload)
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		// handle error
		fmt.Errorf("请求错误：%v\n", err)
	}
	//_ = PrintJson(res)
	return res, nil
}

func (c *Client) InstAssoci(method string, url string, body map[string]interface{}) (map[string]interface{}, error) {
	ms, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer([]byte(ms))
	url = c.BaseUrl + url
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetInstanceList(objId string, body map[string]interface{}) (*interface{}, error) {
	ms, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer([]byte(ms))
	url := c.BaseUrl + "/api/v3/find/instassociation/object/" + objId
	req, err := http.NewRequest("POST", url, payload)
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	if res["bk_error_msg"] != "success" {
		return nil, errors.New("获取pod instance list 错误")
	}
	d := res["data"].(map[string]interface{})["info"]

	return &d, nil
}

//Operator: 取值为：$regex $eq $ne
func (c *Client) GetInstance(objId string, body *Condition)(map[string]interface{}, error){
	ms, err := json.Marshal(*body)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer([]byte(ms))
	url := c.BaseUrl + "/api/v3/find/instassociation/object/" + objId
	req, err := http.NewRequest("POST", url, payload)
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (c *Client) DelInstance(objId string, instId string)(map[string]interface{}, error) {

	url := c.BaseUrl + "/api/v3/delete/instance/object/" + objId + "/inst/"+ instId
	fmt.Printf("url: %s\n", url)
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (c *Client) DelInstancesArray(objId string, body *DelInstances)(map[string]interface{}, error){
	ms, err := json.Marshal(*body)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer([]byte(ms))
	url := c.BaseUrl + "/api/v3/deletemany/instance/object/" + objId
	req, err := http.NewRequest("DELETE", url, payload)

	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//func (c *Client) GetInstance()

func (c *Client) GetAssociList(id string) (map[string]interface{}, error) {
	body := map[string]interface{}{"bk_obj_asst_id": id}
	ms, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer([]byte(ms))
	url := c.BaseUrl + "/api/v3/find/instassociation"
	req, err := http.NewRequest("POST", url, payload)
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (c *Client) DelAssoci(id string) (map[string]interface{}, error) {
	url := c.BaseUrl + "/api/v3/delete/instassociation/" + id
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//获取所有模型内容
func (c *Client) GetModels() (map[string]interface{}, error) {
	url := c.BaseUrl + "/api/v3/find/classificationobject"
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//获取具体某个模型的所有关联
func (c *Client) GetObjAssociation(name string) (map[string]interface{},map[string]interface{}, error) {
	body1 := map[string]interface{}{"condition": map[string]string{"bk_asst_obj_id": name}}
	body2 := map[string]interface{}{"condition": map[string]string{"bk_obj_id": name}}
	ms1, err := json.Marshal(body1)
	ms2, err := json.Marshal(body2)
	if err != nil {
		return nil,nil, err
	}
	payload1 := bytes.NewBuffer([]byte(ms1))
	payload2 := bytes.NewBuffer([]byte(ms2))
	url := c.BaseUrl + "/api/v3/find/objectassociation"
	req1, err := http.NewRequest("POST", url, payload1)
	req2, err := http.NewRequest("POST", url, payload2)
	req1.Header.Set("Content-Type", c.ContentType)
	req1.Header.Set("Cookie", c.CookieStr)
	req2.Header.Set("Content-Type", c.ContentType)
	req2.Header.Set("Cookie", c.CookieStr)
	resp1, err := c.HttpClient.Do(req1)
	resp2, err := c.HttpClient.Do(req2)
	if err != nil {
		return nil,nil, err
	}
	defer resp1.Body.Close()
	defer resp2.Body.Close()
	res1, err := ParseResponse(resp1)
	res2, err := ParseResponse(resp2)
	if err != nil {
		return nil,nil, err
	}
	return res1, res2, nil
}

func (c *Client) ClearAllAssociations() error{
	//step1:获取Objects
	objects, err := c.GetModels()
	if err != nil {
		return err
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
			associRes1, associRes2, err := c.GetObjAssociation(objId)
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
					res, err := c.GetAssociList(bkAsstObjId)
					if err != nil || res["bk_error_msg"] != "success"{
						klog.Errorf("GetAssociList error: %v\n", err)
					}
					for _, item := range res["data"].([]interface{}) {
						id := strconv.Itoa(int(item.(map[string]interface{})["id"].(float64)))
						_, err := c.DelAssoci(id)
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
					res, err := c.GetAssociList(bkAsstObjId)
					if err != nil {
						klog.Errorf("GetAssociList error: %v\n", err)
					}
					for _, item := range res["data"].([]interface{}) {
						id := strconv.Itoa(int(item.(map[string]interface{})["id"].(float64)))
						klog.Infof("id: %s\n", id)
						res, err := c.DelAssoci(id)
						klog.Infof("del result:%s\n", res["bk_error_msg"])
						if err != nil {
							klog.Errorf("DElAssociList error: %v\n", err)
						}
					}
				}
			}
		}
	}
	return nil
}

func (c *Client) ClearAllInstance() error {
	//step1:获取Objects
	objects, err := c.GetModels()
	if err != nil {
		return err
	}
	for _, value := range objects["data"].([]interface{}) {
		//此id为模型分组id
		classificationId := value.(map[string]interface{})["bk_classification_id"].(string)
		if classificationId == "bk_host_manage" || classificationId == "bk_biz_topo" || classificationId == "bk_organization" || classificationId == "bk_network" {
			continue
		}
		var instIds []int
		//获取分组id下的object的数据
		objects := value.(map[string]interface{})["bk_objects"].([]interface{})
		for _, item := range objects {
			objId := item.(map[string]interface{})["bk_obj_id"].(string)
			res, err := c.GetInstanceList(objId, nil)
			if err != nil {
				return err
			}

			for _, item := range (*res).([]interface{}) {
				bkInstId := int(item.(map[string]interface{})["bk_inst_id"].(float64))
				instIds = append(instIds, bkInstId)

			}
			if len(instIds) == 0 {
				continue
			}
			delInsts := DelInstances{Delete: InstIds{Instids: instIds}}
			_, err1 := c.DelInstancesArray(objId, &delInsts)
			if err1 != nil{
				return err1
			}
		}
	}
	return nil
}

func ParseResponse(response *http.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}

	return result, err
}


//解析(map[string]interface{})数据格式并打印出数据
func PrintJson(m map[string]interface{}) map[string]string {
	obj := map[string]string{}
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			if k == "bk_classification_id" {
				obj["bk_classification_id"] = vv
			} else if k == "bk_obj_name" {
				obj["bk_obj_name"] = vv
			} else if k == "bk_obj_id" {
				obj["bk_obj_id"] = vv
			}
			fmt.Println(k, "is string", vv)
		case float64:
			if k == "id" {
				obj["id"] = strconv.FormatFloat(vv, 'E', -1, 64)
			}
			fmt.Println(k, "is float", int64(vv))
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case nil:
			fmt.Println(k, "is nil", "null")
		case map[string]interface{}:
			fmt.Println(k, "is an map:")
			PrintJson(vv)
		default:
			fmt.Println(k, "is of a type I don't know how to handle ", fmt.Sprintf("%T", v))
		}
	}
	return obj
}
