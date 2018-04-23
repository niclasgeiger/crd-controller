package controller

import (
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "test-controller"

type Controller struct {
	podsSynched cache.InformerSynced
	queue       workqueue.RateLimitingInterface
}
