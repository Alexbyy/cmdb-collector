package agent


type Pods struct{
	Name string `json:"bk_inst_name""`
	Namespace string `json:"icp_pod_namespace,omitempty"`
	HostName string `json:"icp_pod_hostname,omitempty"`
	NodeName string `json:"icp_pod_nodename,omitempty"`
	Labels  string `json:"icp_pod_labels,omitempty"`
	Kind string `json:"icp_pod_kind,omitempty"`
	ClusterName string `json:"icp_pod_clustername,omitempty"`
	Status string`json:"icp_pod_status,omitempty"`
	PodIP string `json:"icp_pod_ip,omitempty"`

}

type Container struct {
	Name       string          `json:"name" protobuf:"bytes,1,opt,name=name"`
	Image      string          `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	Command    []string        `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	Args       []string        `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
	WorkingDir string          `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
	Ports      []ContainerPort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,6,rep,name=ports"`
}

type ContainerPort struct {
	Name          string   `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	HostPort      int32    `json:"hostPort,omitempty" protobuf:"varint,2,opt,name=hostPort"`
	ContainerPort int32    `json:"containerPort" protobuf:"varint,3,opt,name=containerPort"`
	Protocol      string `json:"protocol,omitempty" protobuf:"bytes,4,opt,name=protocol,casttype=Protocol"`
	HostIP        string   `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
}

//type PodStatus struct {
//	Phase                      PodPhase          `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=PodPhase"`
//	Conditions                 []PodCondition    `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
//	Message                    string            `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
//	Reason                     string            `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
//	NominatedNodeName          string            `json:"nominatedNodeName,omitempty" protobuf:"bytes,11,opt,name=nominatedNodeName"`
//	HostIP                     string            `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
//	PodIP                      string            `json:"podIP,omitempty" protobuf:"bytes,6,opt,name=podIP"`
//	PodIPs                     []PodIP           `json:"podIPs,omitempty" protobuf:"bytes,12,rep,name=podIPs" patchStrategy:"merge" patchMergeKey:"ip"`
//	StartTime                  *v1.Time      `json:"startTime,omitempty" protobuf:"bytes,7,opt,name=startTime"`
//	InitContainerStatuses      []ContainerStatus `json:"initContainerStatuses,omitempty" protobuf:"bytes,10,rep,name=initContainerStatuses"`
//	ContainerStatuses          []ContainerStatus `json:"containerStatuses,omitempty" protobuf:"bytes,8,rep,name=containerStatuses"`
//	QOSClass                   PodQOSClass       `json:"qosClass,omitempty" protobuf:"bytes,9,rep,name=qosClass"`
//	EphemeralContainerStatuses []ContainerStatus `json:"ephemeralContainerStatuses,omitempty" protobuf:"bytes,13,rep,name=ephemeralContainerStatuses"`
//}

