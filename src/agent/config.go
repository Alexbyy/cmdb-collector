package agent

var Association = map[string][]string{
	"pod_group_container": {"name"},
	"node_run_pod":        {"name"},
	"sts_create_pod":      {},
	"job_create_pod":      {},
	"deploy_create_pod":   {},
	"ds_create_pod":       {"name", "name2"},
}
