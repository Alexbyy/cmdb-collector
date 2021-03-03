package collector

import (
	a "cmdb-collector/src/agent"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

func PreparePodData(item v1.Pod) a.Pods {
	var labelsStr string
	for key := range item.Labels {
		labelsStr += key
	}

	//获取pod的workload:deployment\sts\ds etc
	var kind string
	if len(item.ObjectMeta.OwnerReferences) > 0 {
		kind = item.ObjectMeta.OwnerReferences[0].Kind
	}
	pod := a.Pods{
		Name:        item.Name,
		Namespace:   item.Namespace,
		NodeName:    item.Spec.NodeName,
		HostName:    item.Spec.Hostname,
		ClusterName: item.ClusterName,
		Kind:        kind,
		Labels:      labelsStr,
		Status:      string(item.Status.Phase),
		PodIP:       item.Status.PodIP,
	}

	return pod
}

func PrepareContainerData(item v1.Pod) *[]a.Container {
	containers := item.Spec.Containers
	var res []a.Container
	for _, c := range containers {

		command := strings.Join(c.Command, " ")
		arg := strings.Join(c.Args, " ")
		var ports string
		for _, cp := range c.Ports {
			ports = strconv.Itoa(int(cp.ContainerPort)) + " "
		}

		container := a.Container{
			Name:          c.Name + "_" + item.Name,
			ContainerName: c.Name,
			Image:         c.Image,
			Command:       command,
			Args:          arg,
			WorkingDir:    c.WorkingDir,
			Ports:         ports,
		}
		res = append(res, container)
	}
	return &res
}

func PrepareNodeData(item v1.Node) (*a.Node, error) {
	np := string(item.Status.Phase)
	labels := ""
	for key, val := range item.Labels {
		labels = labels + key + ":" + val + ";"
	}

	node := a.Node{
		Name:        item.Name,
		NodePhase:   np,
		Labels:      labels,
		ClusterName: item.ClusterName,
	}
	return &node, nil
}
