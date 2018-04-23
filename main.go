package main

import (
	"flag"
	"fmt"

	"net/http"

	"github.com/golang/glog"
	niclasgeigerclient "github.com/niclasgeiger/crd-controller/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	port = ":80"
)

var (
	kuberconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master      = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
)
var client *niclasgeigerclient.Clientset

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kuberconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}

	client, err = niclasgeigerclient.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %v", err)
	}
	http.HandleFunc("/", Handle)
	http.ListenAndServe(port, nil)
}

func Handle(writer http.ResponseWriter, request *http.Request) {
	list, err := client.NiclasgeigerV1().Foos("default").List(metav1.ListOptions{})
	if err != nil {
		glog.Fatalf("Error listing all Foos: %v", err)
	}

	foos := make([]string, 0)
	for _, db := range list.Items {
		foos = append(foos, fmt.Sprintf("Foo %s with Bar %q\n", db.Name, db.Spec.Bar))
	}
	output := fmt.Sprintf("%s", foos)
	writer.Write([]byte(output))
}
