package main

import (
	"flag"

	"path/filepath"

	"github.com/niclasgeiger/crd-controller/pkg/controller"
	"github.com/nikhita/custom-database-controller/pkg/signals"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	logrus.Info("Starting Listening")
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := KubeConfig()
	if err != nil {
		logrus.Fatalf("Error getting kubeconfig: %s", err.Error())
	}

	crdController, err := controller.NewController(cfg)
	if err != nil {
		logrus.Fatalf("Error creating crd controller: %s", err.Error())
	}

	if err = crdController.Run(stopCh); err != nil {
		logrus.Fatalf("Error running crdController: %s", err.Error())
	}
}

func KubeConfig() (*rest.Config, error) {
	var kubeconfig *string
	if kubeconfig, err := rest.InClusterConfig(); err == nil {
		return kubeconfig, nil
	}

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	return clientcmd.BuildConfigFromFlags("", *kubeconfig)
}
