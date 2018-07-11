package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/glog"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorln(err)
	}

	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"pods",
		v1.NamespaceAll,
		fields.Everything(),
	)
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				podEvent("add", obj)
			},
			DeleteFunc: func(obj interface{}) {
				podEvent("delete", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				podEvent("update", newObj)
			},
		},
	)

	cache.NewSharedIndexInformer(watchlist,
		&v1.Service{}, 0, cache.Indexers{
			cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
		})
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

func podEvent(event string, obj interface{}) {
	switch event {
	case "add":
		fmt.Printf("pod added: %s \n", obj)
	case "delete":
		fmt.Printf("pod deleted: %s \n", obj)
	case "update":
		fmt.Printf("pod changed \n")
	default:
		fmt.Println("Invalid pod event")
	}
}

//func postToDiscord()
