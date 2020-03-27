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
	"context"
	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog"

	"github.com/duyanghao/velero-volume-controller/cmd/controller/velerovolume"
	"github.com/duyanghao/velero-volume-controller/cmd/controller/velerovolume/config"
	"github.com/duyanghao/velero-volume-controller/pkg/signals"
	"github.com/google/uuid"
)

var (
	argConfigPath string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// load configuration
	cfg, err := config.LoadConfig(argConfigPath)
	if err != nil {
		klog.Fatalf("Error loading veleroVolumeConfig: %s", err.Error())
	}

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	clientCfg, err := clientcmd.BuildConfigFromFlags(cfg.ClusterServerCfg.MasterURL, cfg.ClusterServerCfg.KubeConfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(clientCfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	run := func() {
		kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, 0)

		controller := velerovolume.NewController(cfg.VeleroVolumeCfg, kubeClient, kubeInformerFactory.Core().V1().Pods())

		// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
		// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
		kubeInformerFactory.Start(stopCh)

		if err = controller.Run(1, stopCh); err != nil {
			klog.Fatalf("Error running controller: %s", err.Error())
		}
	}

	// construct lock identity: os.hostname + '_' + uuid.New().String()
	identity, err := os.Hostname()
	if err != nil {
		klog.Fatalf("Error get hostname: %s", err.Error())
	}
	identity = identity + "_" + uuid.New().String()

	// leader election uses the Kubernetes API by writing to a
	// lock object, which can be a LeaseLock object (preferred),
	// a ConfigMap, or an Endpoints (deprecated) object.
	// Conflicting writes are detected and each client handles those actions
	// independently.

	// we use the Lease lock type since edits to Leases are less common
	// and fewer objects in the cluster watch "all Leases".
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      cfg.ClusterServerCfg.LeaseLockName,
			Namespace: cfg.ClusterServerCfg.LeaseLockNamespace,
		},
		Client: kubeClient.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: identity,
		},
	}

	// start the leader election code loop
	leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
				run()
			},
			OnStoppedLeading: func() {
				// we can do cleanup here
				klog.Infof("Leader lost: %s", identity)
				os.Exit(0)
			},
			OnNewLeader: func(id string) {
				// we're notified when new leader elected
				if identity == id {
					// I just got the lock
					return
				}
				klog.Infof("New leader elected: %s", id)
			},
		},
	})
}

func init() {
	flag.StringVar(&argConfigPath, "c", "/cluster-coredns-controller/examples/config.yml", "The configuration filepath for cluster-coredns-controller.")
}
