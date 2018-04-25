package controller

import (
	"time"

	"fmt"

	informers "github.com/niclasgeiger/crd-controller/pkg/client/informers/externalversions"
	lister "github.com/niclasgeiger/crd-controller/pkg/client/listers/niclasgeiger.com/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "test-controller"

type Controller struct {
	podsSynched cache.InformerSynced
	workqueue   workqueue.RateLimitingInterface
	fooLister   lister.FooLister
	k8sClient   *kubernetes.Clientset
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

func NewController(kubeClient *kubernetes.Clientset, factory informers.SharedInformerFactory) *Controller {
	fooInformer := factory.Niclasgeiger().V1().Foos()

	controller := &Controller{
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Foos"),
		fooLister: fooInformer.Lister(),
		k8sClient: kubeClient,
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
func (c *Controller) Run(stopCh <-chan struct{}) error {
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

	logrus.Info("Starting worker")
	go wait.Until(c.runWorker, time.Second, stopCh)

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
	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		logrus.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Foo resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}
	logrus.WithFields(logrus.Fields{
		"namespace": namespace,
		"name":      name,
	}).Info("syncing Handlers")

	// Get the Foo resource with this namespace/name
	foo, err := c.fooLister.Foos(namespace).Get(name)
	if err != nil {
		// The Foo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("foo '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	srv := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      foo.Name,
			Namespace: foo.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,

			Selector: map[string]string{
				"app": "test",
			},
			Ports: []v1.ServicePort{
				{
					Name:     "tcp",
					Protocol: v1.ProtocolTCP,
					Port:     foo.Spec.Port,
					NodePort: foo.Spec.NodePort,
				},
			},
		},
	}
	if _, err = c.k8sClient.CoreV1().Services("default").Create(srv); err != nil {
		logrus.Warn("Could not create Service")
		return err
	}

	return nil
}
