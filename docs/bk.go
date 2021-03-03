package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Model struct {
	Bk_classification_id string `json:"bk_classification_id"`
	Bk_obj_icon          string `json:"bk_obj_icon"`
	Bk_obj_id            string `json:"bk_obj_id"`
	Bk_obj_name          string `json:"bk_obj_name"`
	Bk_supplier_account  string `json:"bk_supplier_account"`
	UserName             string `json:"userName"`
}

type Attr struct {
	Bk_obj_id           string `json:"bk_obj_id"`
	Bk_property_group   string `json:"bk_property_group"`
	Bk_property_id      string `json:"bk_property_id"`
	Bk_property_index   string `json:"bk_property_index"`
	Bk_property_name    string `json:"bk_property_name"`
	Bk_property_type    string `json:"bk_property_type"`
	Bk_supplier_account string `json:"bk_supplier_account"`
	Creator             string `json:"creator"`
	Editable            string `json:"editable"`
	Isrequired          string `json:"isrequired"`
	Option              string `json:"option"`
	Placeholder         string `json:"placeholder"`
	Unit                string `json:"unit"`
}
type AssociObj struct {
	Bk_asst_id       string `json:"bk_asst_id"`
	Bk_asst_obj_id   string `json:"bk_asst_obj_id"`
	Bk_obj_asst_id   string `json:"bk_obj_asst_id"`
	Bk_obj_asst_name string `json:"bk_obj_asst_name"`
	Bk_obj_id        string `json:"bk_obj_id"`
	Mapping          string `json:"mapping"`
}

type Client struct {
	BaseUrl     string
	CookieMap   map[string]string
	CookieStr   string
	ContentType string
	HttpClient  *http.Client
	objs        []map[string]string
}

func newClient(url string) *Client {
	cookies := map[string]string{
		"HTTP_BLUEKING_SUPPLIER_ID": "0",
		"http_scheme":               "http",
		"cc3":                       "MTYxMTkxMzIzMHxOd3dBTkVGQ1NVTkZTVTVHVTA0MVZUTk5XVkUxVlVSRlVFNUVTRmRMVXpWRFZsSTFTRFpHVGtVeldEZFBRMFEyUjFCRlRVOUpVRkU9fBNp_HYCb_mz_B0U210DQL9ZLcp48P2rA1PPJb3CLwZJ",
	}
	cookieStr := "HTTP_BLUEKING_SUPPLIER_ID=0; http_scheme=http; cc3=MTYxMTkxMzIzMHxOd3dBTkVGQ1NVTkZTVTVHVTA0MVZUTk5XVkUxVlVSRlVFNUVTRmRMVXpWRFZsSTFTRFpHVGtVeldEZFBRMFEyUjFCRlRVOUpVRkU9fBNp_HYCb_mz_B0U210DQL9ZLcp48P2rA1PPJb3CLwZJ"
	newClient := &Client{
		BaseUrl:     "http://10.110.19.61:32033",
		ContentType: "application/json;charset=UTF-8",
		CookieMap:   cookies,
		CookieStr:   cookieStr,
		HttpClient:  &http.Client{},
	}
	return newClient
}
func (c *Client) addModel(method string, url string, body Model) {
	ms, err := json.Marshal(body)
	if err != nil {
		fmt.Errorf("json 编译错误：v%\n", err)
	}
	payload := bytes.NewBuffer([]byte(ms))
	url = c.BaseUrl + url
	reqAddModel, err := http.NewRequest(method, url, payload)
	reqAddModel.Header.Set("Content-Type", c.ContentType)
	reqAddModel.Header.Set("Cookie", c.CookieStr)

	resp, err := c.HttpClient.Do(reqAddModel)
	defer resp.Body.Close()
	res, err := ParseResponse(resp)
	if err != nil {
		// handle error
		fmt.Errorf("请求错误：%v\n", err)
	}
	obj := PrintJson(res)
	_ = append(c.objs, obj)
	fmt.Printf("add model res:%v\n", res)
}

func (c *Client) delModel(method string, url string, id string) {
	sign := false
	for _, obj := range c.objs {
		if obj["id"] == id {
			sign = true
			break
		}
	}
	if !sign {
		fmt.Errorf("不存在id为%v的模型", id)
	}

	url = c.BaseUrl + url + id
	body := new([]byte)
	payload := bytes.NewBuffer(*body)
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
	_ = PrintJson(res)
}

func (c *Client) addAttr(method string, url string, body Attr) {
	ms, err := json.Marshal(body)
	if err != nil {
		fmt.Errorf("json 编译错误：v%\n", err)
	}
	payload := bytes.NewBuffer([]byte(ms))
	url = c.BaseUrl + url
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
	_ = PrintJson(res)
}

func (c *Client) objAssociation(method string, url string, body AssociObj) {
	ms, err := json.Marshal(body)
	if err != nil {
		fmt.Errorf("json 编译错误：v%\n", err)
	}
	payload := bytes.NewBuffer([]byte(ms))
	url = c.BaseUrl + url
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
	_ = PrintJson(res)
}
func (c *Client) addInstance(method string, url string, body map[string]string) {
	ms, err := json.Marshal(body)
	if err != nil {
		fmt.Errorf("json 编译错误：v%\n", err)
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
	_ = PrintJson(res)
}
func (c *Client) instAssoci(method string, url string, body map[string]interface{}) {
	ms, err := json.Marshal(body)
	if err != nil {
		fmt.Errorf("json 编译错误：v%\n", err)
	}
	payload := bytes.NewBuffer([]byte(ms))
	url = c.BaseUrl + url
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
	_ = PrintJson(res)
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
