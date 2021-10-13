package main

import (
	"fmt"
	"time"

	mathinformers "Kubewatch/pkg/client/informers/externalversions/myresource/v1alpha1"

	//apiv1Alphav1 "Kubewatch/pkg/apis/myresource/v1alpha1"

	samplescheme "Kubewatch/pkg/client/clientset/versioned/scheme"

	clientset "Kubewatch/pkg/client/clientset/versioned"
	"Kubewatch/pkg/client/clientset/versioned/scheme"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

func NewController(sampleclientset clientset.Interface, queue workqueue.RateLimitingInterface, exampleInformer mathinformers.MyresourceInformer) *Controller {

	runtime.Must(samplescheme.AddToScheme(scheme.Scheme))

	exampleInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old, new interface{}) {

			/*newDepl := new.(*apiv1Alphav1.Myresource)

			oldDepl := old.(*apiv1Alphav1.Myresource)

			if newDepl.ResourceVersion != oldDepl.ResourceVersion {

				key, err := cache.MetaNamespaceKeyFunc(new)
				if err == nil {
					queue.Add(key)
				}

			} else {

				key, err := cache.MetaNamespaceKeyFunc(old)
				if err == nil {
					queue.Add(key)
				}

			}*/

			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}

		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	return &Controller{
		sampleclientset: sampleclientset,
		MathLister:      exampleInformer.Lister(),
		MathSync:        exampleInformer.Informer().HasSynced,
		queue:           queue,
	}
}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	err := c.syncToStdout(key.(string))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)
	return true
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the pod to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *Controller) syncToStdout(key string) error {

	namespace, name, err := cache.SplitMetaNamespaceKey(key)

	math, err := c.MathLister.Myresources(namespace).Get(name)

	if err != nil {
		klog.Errorf("Fetching CRD  with key %s from store failed with %v", key, err)
		return err
	}

	if math.Spec.Operation != "" {

		switch math.Spec.Operation {

		case ("sum"):
			{
				fmt.Printf("Operation SUM   value %d \n", *math.Spec.FirstNum+*math.Spec.SecondNum)

			}

		case ("sub"):
			{
				fmt.Printf("Operation sub   value %d \n", *math.Spec.FirstNum-*math.Spec.SecondNum)

			}
		}

	} else {

		klog.Errorf("Fetching objectmath.Spec.Operation with  key %s from store failed with %v", key, err)
		return err

	}

	return nil
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	klog.Infof("Dropping pod %q out of the queue: %v", key, err)
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()
	klog.Info("Starting Pod controller")

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.MathSync) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Stopping Pod controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}
