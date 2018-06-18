package controller

import (
	"time"

	"fmt"

	niclasgeigerv1 "github.com/niclasgeiger/crd-controller/pkg/apis/niclasgeiger.com/v1"
	informers "github.com/niclasgeiger/crd-controller/pkg/client/informers/externalversions"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	workqueue workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
}

type Action string

const (
	DeleteAction Action = "delete"
	AddAction    Action = "add"
)

type Job struct {
	ID        string
	Action    Action
	Namespace string
	Name      string
}

func (c Controller) EnqueueItem(action Action, obj interface{}) {
	logrus.Info("Added new Foo Object")
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	job := Job{
		ID:        uuid.New(),
		Action:    action,
		Namespace: namespace,
		Name:      name,
	}
	logrus.Info("Adding job to workqueue")
	c.workqueue.AddRateLimited(job)
}

func NewController(restConfig *rest.Config, factory informers.SharedInformerFactory, kubeClient *kubernetes.Clientset) *Controller {

	apiextcsClient, err := apiextcs.NewForConfig(restConfig)
	if err != nil {
		panic(err.Error())
	}
	if err := niclasgeigerv1.CreateCRD(apiextcsClient); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Could not add CRD to K8s API")
		return nil
	}
	userInformer := factory.Niclasgeiger().V1().Users()

	controller := &Controller{
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Users"),
		informer:  userInformer.Informer(),
	}
	// add event listener for CRD
	userInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			controller.EnqueueItem(AddAction, obj)
		},
		DeleteFunc: func(obj interface{}) {
			controller.EnqueueItem(DeleteAction, obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		err := fmt.Errorf("Timed out waiting for caches to sync")
		runtime.HandleError(err)
		return err
	}

	go wait.Until(c.runWorker, time.Second, stopCh)

	// stop Channel to gracefully shut down
	<-stopCh

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}
	err := func(obj interface{}) error {
		// processing of the item is done after exiting the function
		defer c.workqueue.Done(obj)
		var job Job
		var ok bool
		if job, ok = obj.(Job); !ok {
			// Item can't be processed by this queue
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the Job
		if err := c.syncHandler(job); err != nil {
			return fmt.Errorf("error syncing '%s': %s", job, err.Error())
		}
		// Item can be forgotten as it was successfully processed
		c.workqueue.Forget(obj)
		logrus.Infof("Successfully synced '%s'", job)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(job Job) error {
	// Convert the namespace/name string into a distinct namespace and name
	logrus.WithFields(logrus.Fields{
		"namespace": job.Namespace,
		"name":      job.Name,
	}).Info("syncing Handlers")
	switch job.Action {
	case AddAction:
		fmt.Println("add action executed")
	case DeleteAction:
		fmt.Println("delete action executed")
	}

	return fmt.Errorf("unsupported Job Action")
}
