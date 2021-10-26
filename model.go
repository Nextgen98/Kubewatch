package main

import (
	mathlisters "Kubewatch/pkg/client/listers/myresource/v1alpha1"

	appslisters "k8s.io/client-go/listers/apps/v1"

	clientset "Kubewatch/pkg/client/clientset/versioned"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const GroupName = "nextgengroup.com"

type Controller struct {

	// kubeclientset is a standard kubernetes clientset

	sampleclientset clientset.Interface

	Deploymentlisters appslisters.DeploymentLister

	DeploymentSync cache.InformerSynced

	queue workqueue.RateLimitingInterface

	MathLister mathlisters.MyresourceLister

	MathSync cache.InformerSynced
}
