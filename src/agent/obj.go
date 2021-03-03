package agent

type Pods struct {
	Name        string `json:"bk_inst_name""`
	Namespace   string `json:"icp_pod_namespace,omitempty"`
	HostName    string `json:"icp_pod_hostname,omitempty"`
	NodeName    string `json:"icp_pod_nodename,omitempty"`
	Labels      string `json:"icp_pod_labels,omitempty"`
	Kind        string `json:"icp_pod_kind,omitempty"`
	ClusterName string `json:"icp_pod_clustername,omitempty"`
	Status      string `json:"icp_pod_status,omitempty"`
	PodIP       string `json:"icp_pod_ip,omitempty"`
}

type Container struct {
	Name          string `json:"bk_inst_name"`
	ContainerName string `json:"icp_container_name"`
	Image         string `json:"icp_container_image,omitempty"`
	Command       string `json:"icp_container_comm,omitempty"`
	Args          string `json:"icp_container_args,omitempty"`
	WorkingDir    string `json:"icp_container_wd,omitempty"`
	Ports         string `json:"icp_container_ports,omitempty"`
}

type Node struct {
	Name        string `json:"bk_inst_name"`
	NodePhase   string `json:"icp_node_phase"`
	Labels      string `json:"icp_node_labels"`
	ClusterName string `json:"icp_node_cn"`
}
