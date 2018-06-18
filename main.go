package main

import (
	"flag"
	"time"

	"path/filepath"

	clientset "github.com/niclasgeiger/crd-controller/pkg/client/clientset/versioned"
	informers "github.com/niclasgeiger/crd-controller/pkg/client/informers/externalversions"
	"github.com/niclasgeiger/crd-controller/pkg/controller"
	"github.com/nikhita/custom-database-controller/pkg/signals"
	"github.com/sirupsen/logrus"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	kubeconfig string
	masterURL  string
)

func main() {
	logrus.Info("Starting Listening")
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		logrus.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logrus.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	exampleClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		logrus.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	crdInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	crdController := controller.NewController(cfg, crdInformerFactory, kubeClient)

	go kubeInformerFactory.Start(stopCh)
	go crdInformerFactory.Start(stopCh)

	if err = crdController.Run(stopCh); err != nil {
		logrus.Fatalf("Error running crdController: %s", err.Error())
	}
}

func init() {

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
