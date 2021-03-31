package agent

type Pods struct {
	Name        string `json:"bk_inst_name"`
	Id          string `json:"id"`
	Namespace   string `json:"namespace,omitempty"`
	HostName    string `json:"hostname,omitempty"`
	NodeName    string `json:"nodename,omitempty"`
	Labels      string `json:"labels,omitempty"`
	ClusterName string `json:"clustername,omitempty"`
	Status      string `json:"status,omitempty"`
	PodIP       string `json:"ip,omitempty"`
	OwnerReferencesName string `json:"orn"`
	OwnerReferencesType string `json:"ort"`
	//OrnWithNS   string  `json:"orn_ns"`
	OrnId       string  `json:"orn_id"`
}

type Container struct {
	Name          string `json:"bk_inst_name"`
	Id            string `json:"id"`
	ContainerName string `json:"name"`
	PodNameWithNsK8s       string `json:"pod_ns_k8s"`
	ClusterName string `json:"clustername,omitempty"`
	Image         string `json:"image,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	//Command       string `json:"icp_container_comm,omitempty"`
	//Args          string `json:"icp_container_args,omitempty"`
	WorkingDir string `json:"wd,omitempty"`
	Ports      string `json:"ports,omitempty"`
	K8sName string `json:"k8s_name"`

}

type Node struct {
	Name        string `json:"bk_inst_name"`
	Id          string `json:"id"`
	IP          string `json:"ip"`
	NodePhase   string `json:"phase"`
	Labels      string `json:"labels"`
	ClusterName string `json:"cn"`
	K8sName string `json:"k8s_name"`
	FullName string `json:"node_k8s"`
}

type Statefulsets struct {
	Name            string `json:"bk_inst_name"`
	Id              string `json:"id"`
	Namespace       string `json:"ns"`
	ServiceName     string `json:"sn"`
	Replicas        int32  `json:"rp"`
	//ReadyReplicas   int32  `json:"icp_sts_rrp"`
	//CurrentReplicas int32  `json:"icp_sts_crp"`
	//UpdatedReplicas int32  `json:"icp_sts_urp"`
	Selector        string	`json:"selector"`
	NameWithNS  string  `json:"name_ns"`
	ClusterName string `json:"cn"`
	Release     string `json:"release"`
}

type Deployments struct {
	Name                string `json:"bk_inst_name"`
	Id                  string `json:"id"`
	Namespace           string `json:"ns"`
	Replicas            int32  `json:"rp"`
	Selector        string	`json:"selector"`
	NameWithNS  string  `json:"name_ns"`
	ClusterName string `json:"cn"`
	Release     string `json:"release"`
	//UpdatedReplicas     int32  `json:"icp_deploy_urp"`
	//ReadyReplicas       int32  `json:"icp_deploy_rrp"`
	//AvailableReplicas   int32  `json:"icp_deploy_arp"`
	//UnavailableReplicas int32  `json:"icp_deploy_unrp"`
}

type DaemonSets struct {
	Name      string `json:"bk_inst_name"`
	Id        string `json:"id"`
	Namespace string `json:"ns"`
	Selector        string	`json:"selector"`
	NameWithNS  string  `json:"name_ns"`
	ClusterName string `json:"cn"`
	Release     string `json:"release"`
}

type ReplicaSet struct {
	Name      string `json:"bk_inst_name"`
	Id        string `json:"id"`
	Namespace string `json:"ns"`
	Replicas            int32  `json:"rp"`
	Selector        string	`json:"selector"`
	OwnerReferencesName string `json:"orn"`
	OwnerReferencesType string `json:"ort"`
	NameWithNS  string  `json:"name_ns"`
	OrnWithNS   string  `json:"orn_ns"`
	OrnId       string  `json:"orn_id"`
	ClusterName string `json:"cn"`
	Release     string `json:"release"`
}

type App struct {
	Name string `json:"bk_inst_name"`
	ReleaseName string `json:"release_name"`
	NameSpace   string `json:"ns"`
	K8sName string `json:"k8s_name"`
	AppGroup string `json:"app_group"`
}

type AppGroup struct {
	Name string `json:"bk_inst_name"`
}
