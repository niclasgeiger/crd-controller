package main

import (
	"flag"
	"fmt"

	"time"

	"github.com/golang/glog"
	niclasgeigerclient "github.com/niclasgeiger/crd-controller/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kuberconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master      = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
)

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kuberconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}

	client, err := niclasgeigerclient.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %v", err)
	}

	list, err := client.ExampleV1().Foos("default").List(metav1.ListOptions{})
	if err != nil {
		glog.Fatalf("Error listing all Foos: %v", err)
	}

	for _, db := range list.Items {
		fmt.Printf("Foo %s with Bar %q\n", db.Name, db.Spec.Bar)
	}
	time.Sleep(30 * time.Second)
}
