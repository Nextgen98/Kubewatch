package main

import (
	"flag"
	"path/filepath"
	"time"

	clientset "Kubewatch/pkg/client/clientset/versioned"

	appsinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	mathInformer "Kubewatch/pkg/client/informers/externalversions"

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

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.

	kClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	exampleClient, err := clientset.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	//CRD informer factory
	exampleInformerFactory := mathInformer.NewSharedInformerFactory(exampleClient, time.Second*30)

	//app informer factory
	deployemntInformerFactory := appsinformers.NewSharedInformerFactory(kClient, time.Second*30)

	controller := NewController(queue, exampleInformerFactory.Math().V1alpha1().Myresources(), deployemntInformerFactory.Apps().V1().Deployments(), exampleClient)

	stop := make(chan struct{})

	// Now let's start the Informer

	deployemntInformerFactory.Start(stop)

	exampleInformerFactory.Start(stop)

	// Now let's start the controller

	defer close(stop)
	go controller.Run(2, stop)

	// Wait forever
	select {}
}
