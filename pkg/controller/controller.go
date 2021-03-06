package controller

import (
	"fmt"

	"time"

	"github.com/niclasgeiger/crd-controller/pkg/client/clientset/versioned"
	informers "github.com/niclasgeiger/crd-controller/pkg/client/informers/externalversions"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	informer cache.SharedIndexInformer
	factory  informers.SharedInformerFactory
}

type ResourceEventHandler struct {
}

func (r ResourceEventHandler) OnAdd(obj interface{}) {
	fmt.Println("On Add")
	fmt.Println(obj)
}

func (r ResourceEventHandler) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("On Update")
	fmt.Println("old:")
	fmt.Println(oldObj)
	fmt.Println("new:")
	fmt.Println(newObj)
}

func (r ResourceEventHandler) OnDelete(obj interface{}) {
	fmt.Println("On Delete")
	fmt.Println(obj)
}

func NewController(restConfig *rest.Config) (*Controller, error) {
	crdClient, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	crdInformerFactory := informers.NewSharedInformerFactory(crdClient, time.Second*30)
	crdInformer := crdInformerFactory.Example().V1().Foos().Informer()

	controller := &Controller{
		informer: crdInformer,
		factory:  crdInformerFactory,
	}
	controller.informer.AddEventHandler(new(ResourceEventHandler))

	return controller, nil
}

func (c Controller) Run(stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	// TODO: is this sufficient?
	//go c.informer.Run(stopCh)
	go c.factory.Start(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		err := fmt.Errorf("timed out waiting for caches to sync")
		runtime.HandleError(err)
		return err
	}

	// stop Channel to run until a stop event is receivred
	<-stopCh

	return nil
}
