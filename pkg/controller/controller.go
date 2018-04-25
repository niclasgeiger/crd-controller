package controller

import (
	"time"

	informers "github.com/niclasgeiger/crd-controller/pkg/client/informers/externalversions"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "test-controller"

type Controller struct {
	podsSynched cache.InformerSynced
	workqueue   workqueue.RateLimitingInterface
}

func (c Controller) EnqueueItem(obj interface{}) {
	logrus.Info("Added new Foo Object")
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	logrus.Info("Adding key to workqueue")
	c.workqueue.AddRateLimited(key)
}

func NewController(factory informers.SharedInformerFactory) *Controller {
	controller := &Controller{
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Foos"),
	}

	// add event listener for CRD
	factory.Niclasgeiger().V1().Foos().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.EnqueueItem,
		DeleteFunc: func(obj interface{}) {

		},
		UpdateFunc: func(oldObj, newObj interface{}) {

		},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	logrus.Info("Starting Foo controller")

	// Wait for the caches to be synced before starting workers
	logrus.Info("Waiting for informer caches to sync")
	/*
		if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.foosSynced); !ok {
			return errors.New("failed to wait for caches to sync")
		}
	*/

	logrus.Info("Starting workers")
	// Launch two workers to process Foo resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	logrus.Info("Started workers")
	<-stopCh
	logrus.Info("Shutting down workers")

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
		logrus.Info("shutting down worker")
		return false
	}

	logrus.Infof("New Object: %v\n", obj)

	return true
}
