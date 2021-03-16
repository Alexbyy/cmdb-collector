package agent

var Association = map[string]map[string]string{
	"pod_group_container": map[string]string{"pod": "bk_inst_name", "container": "icp_pod_name"},
	"node_run_pod":        map[string]string{"node": "bk_inst_name", "pod": "icp_pod_nodename"},
	//"sts_create_pod":      {},
	//"job_create_pod":      {},
	//"deploy_create_pod":   {},
	//"ds_create_pod":       {"name", "name2"},
}
