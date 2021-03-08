package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	//fmt.Printf("body %v\n", req.Body)
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

func (c *Client) InstAssoci(method string, url string, body map[string]interface{})(map[string]interface{}, error)  {
	ms, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer([]byte(ms))
	url = c.BaseUrl + url
	req, err := http.NewRequest(method, url, payload)
	if err != nil{
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

func (c *Client) GetInstanceList(url string, body map[string]interface{}) (*interface{}, error) {
	ms, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer([]byte(ms))
	url = c.BaseUrl + url
	req, err := http.NewRequest("POST", url, payload)
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Cookie", c.CookieStr)
	resp, err := c.HttpClient.Do(req)
	if err != nil{
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

func (c *Client) GetAssociList(id string)(map[string]interface{}, error) {
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
	if err != nil{
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
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	return res, nil
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
