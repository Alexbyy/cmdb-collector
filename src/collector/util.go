package collector

import (
	a "cmdb-collector/src/agent"
	app "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"strconv"
)

func PreparePodData(item v1.Pod) *a.Pods {
	var labelsStr string
	for key := range item.Labels {
		labelsStr += key
	}

	//获取pod的workload:deployment\sts\ds etc
	var ownerReferencesName string
	var ownerReferencesType string
	if len(item.OwnerReferences) > 0{
		ownerReferencesName = item.ObjectMeta.OwnerReferences[0].Name
		ownerReferencesType = item.ObjectMeta.OwnerReferences[0].Kind
	}

	pod := a.Pods{
		Name:        item.Name,
		Id:          string(item.UID),
		Namespace:   item.Namespace,
		NodeName:    item.Spec.NodeName,
		HostName:    item.Spec.Hostname,
		ClusterName: item.ClusterName,
		Labels:      labelsStr,
		Status:      string(item.Status.Phase),
		PodIP:       item.Status.PodIP,
		OwnerReferencesName: ownerReferencesName,
		OwnerReferencesType: ownerReferencesType,
	}

	return &pod
}

func PrepareContainerData(item v1.Pod) *[]a.Container {
	containers := item.Spec.Containers
	var res []a.Container
	for _, c := range containers {

		//command := ""
		//command = strings.Join(c.Command, " ")
		//if len(command) > 20{
		//	command = command[0: 20]
		//}
		//
		//arg := ""
		//arg = strings.Join(c.Args, " ")
		//if len(arg) > 20{
		//	arg = arg[0:20]
		//}
		var ports string
		for _, cp := range c.Ports {
			ports = strconv.Itoa(int(cp.ContainerPort)) + " "
		}

		container := a.Container{
			Name:          c.Name,
			Id:            c.Name + "_" + string(item.UID),
			ContainerName: c.Name,
			PodName:       item.Name,
			Image:         c.Image,
			//Command:       command,
			//Args:          arg,
			WorkingDir: c.WorkingDir,
			Ports:      ports,
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
	address := ""
	for _, item := range item.Status.Addresses {
		address = address + string(item.Type) + ":" + item.Address
	}

	node := a.Node{
		Name:        item.Name,
		Id:          string(item.UID),
		NodePhase:   np,
		Labels:      labels,
		ClusterName: item.ClusterName,
		IP:          address,
	}
	return &node, nil
}

func PrepareStsData(item app.StatefulSet) (*a.Statefulsets, error) {
	sel := ""
	for key, val := range item.Spec.Selector.MatchLabels {
		sel = sel + key + ":" + val + ";"
	}
	sts := a.Statefulsets{
		Name:            item.Name,
		Id:              string(item.UID),
		Namespace:       item.Namespace,
		ServiceName:     item.Spec.ServiceName,
		Replicas:        item.Status.Replicas,
		Selector:        sel,

	}
	return &sts, nil
}

func PrepareDeployData(item app.Deployment) (*a.Deployments, error) {
	sel := ""
	for key, val := range item.Spec.Selector.MatchLabels {
		sel = sel + key + ":" + val + ";"
	}
	deploy := a.Deployments{
		Name:                item.Name,
		Id:                  string(item.UID),
		Namespace:           item.Namespace,
		Replicas:            item.Status.Replicas,
		Selector:        sel,

	}
	return &deploy, nil
}

func PrepareDsData(item app.DaemonSet) (*a.DaemonSets, error) {
	sel := ""
	for key, val := range item.Spec.Selector.MatchLabels {
		sel = sel + key + ":" + val + ";"
	}
	ds := a.DaemonSets{
		Name:      item.Name,
		Id:        string(item.UID),
		Namespace: item.Namespace,
		Selector:        sel,
	}
	return &ds, nil
}
