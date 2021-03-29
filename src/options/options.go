package options

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"os"
)

type Options struct {
	ConfigPath string
	LcmBaseUrl string
	LcmBranch string
	LcmSite string
	K8s     *[]interface{}

	flags *pflag.FlagSet
}

// NewOptions returns a new instance of `Options`.
func NewOptions() *Options {
	return &Options{}
}
// AddFlags populated the Options struct from the command line arguments passed.
func (o *Options) AddFlags() {
	o.flags = pflag.NewFlagSet("", pflag.ExitOnError)
	// add klog flags
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	defer klog.Flush()
	o.flags.AddGoFlagSet(klogFlags)
	o.flags.Lookup("logtostderr").Value.Set("true")
	o.flags.Lookup("logtostderr").DefValue = "true"
	o.flags.Lookup("logtostderr").NoOptDefVal = "true"

	o.flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		o.flags.PrintDefaults()
	}


	o.flags.StringVar(&o.ConfigPath, "config_path", "/config.json", "The path of config file")
	o.flags.StringVar(&o.LcmBaseUrl, "lcm_base_url", "http://10.110.18.31:30977", "The base url of lcm")
	o.flags.StringVar(&o.LcmSite, "lcm_site", "icpshiptest", "The site's name of lcm")
	o.flags.StringVar(&o.LcmBranch, "lcm_branch", "master", "The site's branch of lcm")
}

// Parse parses the flag definitions from the argument list.
func (o *Options) Parse() error {
	err := o.flags.Parse(os.Args)
	return err
}

// Usage is the function called when an error occurs while parsing flags.
func (o *Options) Usage() {
	o.flags.Usage()
}