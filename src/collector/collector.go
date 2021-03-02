package collector

import (
	a "cmdb-collector/src/agent"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Collector struct {
	namespace []string
	client *kubernetes.Clientset
}

func NewCollector() *Collector {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return &Collector{
		client: clientset,
	}
}

func(c *Collector) GetPods(ns string) *[]a.Pods{
	var podList []a.Pods
	pods, err := c.client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	//pod, err := c.client.CoreV1().Pods("monitoring").Get(context.TODO(), "prometheus",  metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	
	for _, item := range pods.Items{
		data := PreparePodData(item, ns)
		podList = append(podList, data)
		fmt.Printf("pod item:  %s\n", data)
	}

	return &podList
}


func(c *Collector) GetNamespaces(){
	ns, err := c.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d namespaces in the cluster\n", len(ns.Items))

	for _, item := range ns.Items{
		fmt.Printf("namespace item: name: %s, spec: %s, status: %s\n", item.Name, item.Spec, item.Status)
	}

}

func(c *Collector) GetNodes(){
	nodes, err := c.client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	for _, item := range nodes.Items{
		fmt.Printf("node item: name: %s, spec: %s, status: %s\n", item.Name, item.Spec, item.Status)
	}
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))

}
