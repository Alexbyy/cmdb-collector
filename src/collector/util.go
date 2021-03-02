package collector

import (
	a "cmdb-collector/src/agent"
	v1 "k8s.io/api/core/v1"
)

func PreparePodData(item v1.Pod, ns string) a.Pods {
	var labelsStr string
	for key := range item.Labels{
		labelsStr += key
	}

	//获取pod的workload:deployment\sts\ds etc
	var kind string
	if len(item.ObjectMeta.OwnerReferences) > 0 {
		kind = item.ObjectMeta.OwnerReferences[0].Kind
	}
	pod := a.Pods{
		Name: item.Name,
		Namespace: ns,
		NodeName: item.Spec.NodeName,
		HostName: item.Spec.Hostname,
		ClusterName: item.ClusterName,
		Kind: kind,
		Labels: labelsStr,
		Status: string(item.Status.Phase),
		PodIP: item.Status.PodIP,
	}

	return pod
}

