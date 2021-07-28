package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func getPods(clientset kubernetes.Clientset) {
	// get pods from default namespace
	// you should know api version of resource you want to interact with, in this case for pods its corev1
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}

func getDeployments(clientset kubernetes.Clientset) {
	// get deployments in default namespace
	deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d deployments in the cluster\n", len(deployments.Items))
}

func main() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "listitems"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		fmt.Println("Error loading listitems")
		panic(err.Error())
	}

	// create clientset
	// set of clients from different api versions
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	for {
		getPods(*clientset)
		getDeployments(*clientset)
	}
}
