package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/phineas/kubehook/discord"
	"github.com/phineas/kubehook/webhook"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Config struct {
	Namespace string `json:"kubernetesNamespace"`
	Events    struct {
		PodAdd    bool `json:"podAdd"`
		PodDelete bool `json:"podDelete"`
		PodUpdate bool `json:"podUpdate"`
	} `json:"events"`
	Services struct {
		Discord string `json:"discord"`
	} `json:"services"`
}

type PodCBRef struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Reference  struct {
		Kind            string `json:"kind"`
		Namespace       string `json:"namespace"`
		Name            string `json:"name"`
		UID             string `json:"uid"`
		APIVersion      string `json:"apiVersion"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"reference"`
}

type PodEvent struct {
	Event string
	Obj   interface{}
}

func main() {
	file, _ := os.Open("config.json")
	defer file.Close()
	configuration := Config{}
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("Config loaded!")

	hook := webhook.NewHook(configuration.Services.Discord)

	//UNCOMMENT BELOW FOR LOCAL DEVELOPMENT
	/*kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}*/

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"pods",
		v1.NamespaceDefault,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if configuration.Events.PodAdd {
					podEvent("add", obj, hook)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if configuration.Events.PodDelete {
					podEvent("delete", obj, hook)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				podEvent("update", newObj, hook)
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

func podEvent(event string, obj interface{}, hook webhook.Webhook) {
	pod := obj.(*v1.Pod)
	r := PodCBRef{}
	err := json.Unmarshal([]byte(pod.Annotations["kubernetes.io/created-by"]), &r)
	if err != nil {
		log.Fatal(err)
	}

	nodeID := pod.Spec.NodeName
	if len(nodeID) < 1 {
		nodeID = "Pending Allocation"
	}

	switch event {
	case "add":
		embed := discord.Embed{Color: 51025, Footer: discord.Footer{Text: "Kubernetes", IconURL: "http://www.stickpng.com/assets/images/58480a44cef1014c0b5e4917.png"}, Description: "`" + r.Reference.Name + "` was scaled up\n**New Pod ID>** `" + pod.Name + "`\n**Node>** `" + nodeID + "`" + "\n**Phase>** `" + string(pod.Status.Phase) + "`"}
		postToDiscord(hook, embed)
	case "delete":
		embed := discord.Embed{Color: 15689877, Footer: discord.Footer{Text: "Kubernetes", IconURL: "http://www.stickpng.com/assets/images/58480a44cef1014c0b5e4917.png"}, Description: "`" + r.Reference.Name + "` was scaled down\n**New Pod ID>** `" + pod.Name + "`\n**Node>** `" + pod.Spec.NodeName + "`" + "\n**Phase>** `" + string(pod.Status.Phase) + "`"}
		postToDiscord(hook, embed)
	case "update":
		fmt.Printf("pod changed \n")
	default:
		fmt.Println("Invalid pod event")
	}
}

func postToDiscord(hook webhook.Webhook, embed discord.Embed) {
	s := []discord.Embed{}
	s = append(s, embed)
	msg := discord.Message{Content: "", Embeds: s}
	hook.Post(msg)
}
