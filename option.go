package main

import "flag"

type options struct {
	version    bool
	configFile string
}

func newOptions() *options {
	return &options{}
}

func (o *options) InstallFlags() {
	flag.BoolVar(&o.version, "version", false, "show the app version")
	flag.StringVar(&o.configFile, "config", "../conf/server.yaml", "set config file")

	flag.Parse()
}
