package options

import (
	"flag"
)
var configPath = flag.String("config_path", "/config.json", "The path of config file")
var lcmBaseUrl =  flag.String("lcm_base_url", "http://10.110.18.31:30977", "The base url of lcm")
var lcmSite = flag.String("lcm_site", "icpshiptest", "The site's name of lcm")
var lcmBranch = flag.String("lcm_branch", "xh-test-lma", "The site's branch of lcm")
var cron = flag.String("cron", "0 0 1 * * ?", "the cron param specifies the time when collector should start again")

type Options struct {
	ConfigPath string
	LcmBaseUrl string
	LcmBranch string
	LcmSite string
	Cron string
	K8s     *[]interface{}

}

// NewOptions returns a new instance of `Options`.
func NewOptions() *Options {
	flag.Parse()
	return &Options{
		ConfigPath: *configPath,
		LcmBaseUrl: *lcmBaseUrl,
		LcmBranch:  *lcmBranch,
		LcmSite:    *lcmSite,
		K8s:        nil,
		Cron: *cron,
	}
}


