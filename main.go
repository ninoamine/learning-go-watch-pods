package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
	"fmt"
)

func main() {

	namespace := flag.String("namespace", "default", "kubernetes namespace")
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "path to kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	watcher, err := clientset.CoreV1().Pods(*namespace).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			continue
		}
		fmt.Printf("[%s] Pod %s: %s (%s)\n", event.Type, pod.Name, pod.Status.Phase, pod.Status.PodIP)
	}




	




}