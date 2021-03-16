package agent

type Pods struct {
	Name        string `json:"bk_inst_name"`
	Id          string `json:"icp_pod_id"`
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
	Id            string `json:"icp_container_id"`
	ContainerName string `json:"icp_container_name"`
	PodName       string `json:"icp_pod_name"`
	Image         string `json:"icp_container_image,omitempty"`
	//Command       string `json:"icp_container_comm,omitempty"`
	//Args          string `json:"icp_container_args,omitempty"`
	WorkingDir string `json:"icp_container_wd,omitempty"`
	Ports      string `json:"icp_container_ports,omitempty"`
}

type Node struct {
	Name        string `json:"bk_inst_name"`
	Id          string `json:"icp_node_id"`
	IP          string `json:"icp_node_ip"`
	NodePhase   string `json:"icp_node_phase"`
	Labels      string `json:"icp_node_labels"`
	ClusterName string `json:"icp_node_cn"`
}

type Statefulsets struct {
	Name            string `json:"bk_inst_name"`
	Id              string `json:"icp_sts_id"`
	Namespace       string `json:"icp_sts_ns"`
	ServiceName     string `json:"icp_sts_sn"`
	Replicas        int32  `json:"icp_sts_rp"`
	//ReadyReplicas   int32  `json:"icp_sts_rrp"`
	//CurrentReplicas int32  `json:"icp_sts_crp"`
	//UpdatedReplicas int32  `json:"icp_sts_urp"`
	Selector        string	`json:"icp_sts_selector"`
}

type Deployments struct {
	Name                string `json:"bk_inst_name"`
	Id                  string `json:"icp_deploy_id"`
	Namespace           string `json:"icp_deploy_ns"`
	Replicas            int32  `json:"icp_deploy_rp"`
	Selector        string	`json:"icp_deploy_selector"`
	//UpdatedReplicas     int32  `json:"icp_deploy_urp"`
	//ReadyReplicas       int32  `json:"icp_deploy_rrp"`
	//AvailableReplicas   int32  `json:"icp_deploy_arp"`
	//UnavailableReplicas int32  `json:"icp_deploy_unrp"`
}

type DaemonSets struct {
	Name      string `json:"bk_inst_name"`
	Id        string `json:"icp_ds_id"`
	Namespace string `json:"icp_ds_ns"`
	Selector        string	`json:"icp_ds_selector"`
}
