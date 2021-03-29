package collector

import (
	"cmdb-collector/src/options"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"net/http"
)

type Collector struct {
	namespace []string
	ContentType string
	client    *kubernetes.Clientset
	options *options.Options
	HttpClient  *http.Client

}


func NewCollector(opts *options.Options) *Collector {

	c := &Collector{}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	c.client = clientset
	c.options = opts
	c.GetNamespaces()
	c.HttpClient = &http.Client{}
	if err != nil {
		panic(err.Error())
	}
	return c
}

func (c *Collector) GetObjData(id string) (*[]interface{}, error) {
	var res []interface{}

	if id == "node" {
		r, err := c.GetNodes()
		return r, err
	}
	if id == "pod" {
		fmt.Printf("nm>>>>: %v\n", c.namespace)
		for _, ns := range c.namespace {
			r, err := c.GetPods(ns)
			if err != nil{
				return nil, err
			}
			res = append(res, *r...)
		}

		return &res, nil
	}
	if id == "container" {
		for _, ns := range c.namespace {
			r, err := c.GetContainers(ns)
			if err != nil{
				return nil, err
			}
			res = append(res, *r...)
		}
		return &res, nil

	}
	if id == "deploy" {
		for _, ns := range c.namespace {
			r, err := c.GetDeployments(ns)
			if err != nil{
				return nil, err
			}
			res = append(res, *r...)
		}
		return &res, nil

	}
	if id == "sts" {
		for _, ns := range c.namespace {
			r, err := c.GetStatefulsets(ns)
			if err != nil{
				return nil, err
			}
			res = append(res, *r...)
		}
		return &res, nil
	}
	if id == "ds" {
		for _, ns := range c.namespace {
			r, err := c.GetDaemonSets(ns)
			if err != nil{
				return nil, err
			}
			res = append(res, *r...)
		}
		return &res, nil
	}
	if id == "rc"{
		for _, ns := range c.namespace{
			r, err := c.GetReplicaSets(ns)
			if err != nil{
				return nil, err
			}
			res = append(res, *r...)

		}
		return &res, nil
	}

	return nil, errors.New("未知object id")
}

func (c *Collector) GetPods(ns string) (*[]interface{}, error) {
	var podList []interface{}
	pods, err := c.client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, item := range pods.Items {
		data := PreparePodData(item)
		podList = append(podList, *data)
	}

	return &podList, nil
}

func (c *Collector) GetNodes() (*[]interface{}, error) {
	nodes, err := c.client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	var nodeList []interface{}
	for _, item := range nodes.Items {
		node, err := PrepareNodeData(item)
		if err != nil {
			return nil, err
		}
		nodeList = append(nodeList, *node)
	}
	if err != nil {
		return nil, err
	}
	klog.Infof("There are %d nodes in the cluster\n", len(nodes.Items))
	return &nodeList, nil
}

func (c *Collector) GetContainers(ns string) (*[]interface{}, error) {
	var containers []interface{}
	pods, err := c.client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, item := range pods.Items {
		con := PrepareContainerData(item)
		if len(*con) > 0{
			for _, item := range *con {
				containers = append(containers, item)
			}
		}
	}
	return &containers, nil
}

func (c *Collector) GetNamespaces() {
	ns, err := c.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d namespaces in the cluster\n", len(ns.Items))

	for _, item := range ns.Items {
		//fmt.Printf("namespace item: name: %s, spec: %s, status: %s\n", item.Name, item.Spec, item.Status)
		c.namespace = append(c.namespace,item.ObjectMeta.Name)
	}

}


func (c *Collector) GetStatefulsets(ns string) (*[]interface{}, error) {
	sts, err := c.client.AppsV1().StatefulSets(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var stsList []interface{}
	for _, item := range sts.Items {
		sts, err := PrepareStsData(item)
		if err != nil {
			return nil, err
		}
		stsList = append(stsList, *sts)
	}
	return &stsList, nil
}

func (c *Collector) GetDeployments(ns string) (*[]interface{}, error) {
	deploy, err := c.client.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var deployList []interface{}
	for _, item := range deploy.Items {
		deploy, err := PrepareDeployData(item)
		if err != nil {
			return nil, err
		}
		deployList = append(deployList, *deploy)
	}
	return &deployList, nil
}

func (c *Collector) GetDaemonSets(ns string) (*[]interface{}, error) {
	ds, err := c.client.AppsV1().DaemonSets(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var dsList []interface{}
	for _, item := range ds.Items {
		ds, err := PrepareDsData(item)
		if err != nil {
			return nil, err
		}
		dsList = append(dsList, *ds)
	}
	return &dsList, nil
}

func (c *Collector) GetReplicaSets(ns string) (*[]interface{}, error) {
	rc, err := c.client.AppsV1().ReplicaSets(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var list []interface{}
	for _, item := range rc.Items {
		rc, err := PrepareRCData(item)
		if err != nil {
			return nil, err
		}
		list = append(list, *rc)
	}
	return &list, nil
}

//func (c *Collector) GetK8s()(*[]interface{},error)  {
//
//	url := c.options.LcmBaseUrl +  "/lcm/v1/sites/" + c.options.LcmSite + "/branches/" + c.options.LcmBranch + "/kuberneteses"
//	req, err := http.NewRequest("POST", url, nil)
//	if err != nil {
//		return nil, err
//	}
//	req.Header.Set("Content-Type", c.ContentType)
//	resp, err := c.HttpClient.Do(req)
//	defer resp.Body.Close()
//	res, err := ParseResponse(resp)
//	if err != nil {
//		return nil, err
//	}
//	return res, nil
//
//}

func ParseResponse(response *http.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}

	return result, err
}
