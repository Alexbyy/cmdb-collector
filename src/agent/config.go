package agent



var Association = map[string][]string{
	"pod_group_container": []string{"name"},
	"node_run_pod": []string{"name"},
	"sts_create_pod": []string{},
	"job_create_pod":[]string{},
	"deploy_create_pod": []string{},
	"ds_create_pod": []string{"name","name2"},
}