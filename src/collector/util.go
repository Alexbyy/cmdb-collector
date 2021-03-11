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
	var kind string
	if len(item.ObjectMeta.OwnerReferences) > 0 {
		kind = item.ObjectMeta.OwnerReferences[0].Kind
	}
	pod := a.Pods{
		Name:        item.Name,
		Id:          string(item.UID),
		Namespace:   item.Namespace,
		NodeName:    item.Spec.NodeName,
		HostName:    item.Spec.Hostname,
		ClusterName: item.ClusterName,
		Kind:        kind,
		Labels:      labelsStr,
		Status:      string(item.Status.Phase),
		PodIP:       item.Status.PodIP,
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
			Name:          c.Name + "_" + item.Name,
			Id:            string(item.UID),
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

	node := a.Node{
		Name:        item.Name,
		Id:          string(item.UID),
		NodePhase:   np,
		Labels:      labels,
		ClusterName: item.ClusterName,
	}
	return &node, nil
}

func PrepareStsData(item app.StatefulSet) (*a.Statefulsets, error) {
	sts := a.Statefulsets{
		Name:            item.Name,
		Id:              string(item.UID),
		Namespace:       item.Namespace,
		ServiceName:     item.Spec.ServiceName,
		Replicas:        item.Status.Replicas,
		ReadyReplicas:   item.Status.ReadyReplicas,
		CurrentReplicas: item.Status.CurrentReplicas,
		UpdatedReplicas: item.Status.UpdatedReplicas,
	}
	return &sts, nil
}

func PrepareDeployData(item app.Deployment) (*a.Deployments, error) {
	deploy := a.Deployments{
		Name:                item.Name,
		Id:                  string(item.UID),
		Namespace:           item.Namespace,
		Replicas:            item.Status.Replicas,
		ReadyReplicas:       item.Status.ReadyReplicas,
		UpdatedReplicas:     item.Status.UpdatedReplicas,
		AvailableReplicas:   item.Status.AvailableReplicas,
		UnavailableReplicas: item.Status.UnavailableReplicas,
	}
	return &deploy, nil
}

func PrepareDsData(item app.DaemonSet) (*a.DaemonSets, error) {
	ds := a.DaemonSets{
		Name:      item.Name,
		Id:        string(item.UID),
		Namespace: item.Namespace,
	}
	return &ds, nil
}
