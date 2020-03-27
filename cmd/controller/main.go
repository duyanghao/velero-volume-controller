/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/duyanghao/velero-volume-controller/cmd/controller/velerovolume"
	"github.com/duyanghao/velero-volume-controller/cmd/controller/velerovolume/config"
	"github.com/duyanghao/velero-volume-controller/pkg/signals"
)

var (
	argConfigPath string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// load configuration
	veleroVolumeCfg, err := config.LoadConfig(argConfigPath)
	if err != nil {
		klog.Fatalf("Error loading veleroVolumeConfig: %s", err.Error())
	}

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(veleroVolumeCfg.ClusterServerCfg.MasterURL, veleroVolumeCfg.ClusterServerCfg.KubeConfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, 0)

	controller := velerovolume.NewController(veleroVolumeCfg.VeleroVolumeCfg, kubeClient, kubeInformerFactory.Core().V1().Pods())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)

	if err = controller.Run(1, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&argConfigPath, "c", "/cluster-coredns-controller/examples/config.yml", "The configuration filepath for cluster-coredns-controller.")
}
