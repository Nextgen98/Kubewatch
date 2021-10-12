package main

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	clientset "Kubewatch/pkg/client/clientset/versioned"

	mathInformer "Kubewatch/pkg/client/informers/externalversions"

	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

func main() {
	var kubeconfig *string
	//var master string

	home := homedir.HomeDir()

	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	//flag.StringVar(&kubeconfig, "kubeconfig", "", "") // flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	//  flag.StringVar(&master, "master", "", "")         //flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	//	flag.Parse()

	// creates the connection
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		klog.Fatal(err)
	}

	// creates the clientset
	kclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	ctx := context.Background()

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return kclientset.CoreV1().Pods(meta_v1.NamespaceAll).List(ctx, options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return kclientset.CoreV1().Pods(meta_v1.NamespaceAll).Watch(ctx, options)
			},
		},
		&api_v1.Pod{},
		0, //Skip resync
		cache.Indexers{},
	)

	exampleClient, err := clientset.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	exampleInformerFactory := mathInformer.NewSharedInformerFactory(exampleClient, time.Second*30)

	controller := NewController(kclientset, exampleClient, queue, informer, exampleInformerFactory.Math().V1alpha1().Maths())

	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(2, stop)

	// Wait forever
	select {}
}
