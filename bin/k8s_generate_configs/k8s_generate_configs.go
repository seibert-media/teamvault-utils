package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/bborbe/kubernetes_tools/config"
	"github.com/bborbe/kubernetes_tools/file"
	"github.com/bborbe/kubernetes_tools/model"
	model_generator "github.com/bborbe/kubernetes_tools/model/generator"
	"github.com/golang/glog"
)

const (
	PARAMETER_CONFIG = "config"
)

var (
	configPtr = flag.String(PARAMETER_CONFIG, "", "config json file")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := do()
	if err != nil {
		glog.Exit(err)
	}
}

func do() error {
	configPath := config.ConfigPath(*configPtr)
	if len(configPath) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_CONFIG)
	}
	glog.V(2).Infof("config path: %s", configPath)
	configPath, err := configPath.NormalizePath()
	if err != nil {
		glog.Warningf("normalize path '%s' failed", configPath)
		return err
	}
	glog.V(2).Infof("normalized config path: %s", configPath)
	config, err := configPath.ParseConfig()
	if err != nil {
		glog.Warningf("parse config '%s' failed: %v", config, err)
		return err
	}
	cluster, err := model_generator.Generate(config)
	if err != nil {
		glog.Warningf("generate model failed: %v", err)
		return err
	}
	features := model.Features{
		Kvm: config.Features.Kvm,
	}
	if err := cluster.Validate(features); err != nil {
		glog.Warningf("validate model failed: %v", err)
		return err
	}
	if err := file.Generate(features, *cluster); err != nil {
		glog.Warningf("write configs failed: %v", err)
		return err
	}
	glog.V(2).Infof("generate kubernetes cluster configs completed")
	return nil
}
