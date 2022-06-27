package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	cfg *AppConfig

	namespacestoWatch []string
)

func NewRestConfig() (*rest.Config, error) {
	kubeConfig := os.Getenv("KUBECONFIG")

	var config *rest.Config
	var err error
	if kubeConfig != "" {
		// creates a client from the kube config
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		// creates a in-cluster config
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes config: %w", err)
	}

	return config, nil
}

func main() {
	// Get application configuration
	cfg = NewAppConfig()

	// Create config for the Kubernetes client
	config, err := NewRestConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	// Split namespaces list to slice
	namespacestoWatch = strings.Split(cfg.GetString("NAMESPACES"), ",")

	// Stop signal for the informer
	stopper := make(chan struct{})
	defer close(stopper)

	factory := informers.NewSharedInformerFactory(clientset, 0)
	eventsInformer := factory.Core().V1().Events()
	informer := eventsInformer.Informer()

	defer runtime.HandleCrash()

	log.Printf("Starting %s. Output = %s, Namespaces to watch %v\n",
		AppName,
		cfg.GetString("OUTPUT"),
		namespacestoWatch,
	)

	// Start informer
	go factory.Start(stopper)

	// Start to sync and call list
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	// Register event handlers
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    AddEvent,
		UpdateFunc: UpdateEvent,
		DeleteFunc: DeleteEvent,
	})

	informer.Run(stopper)

	<-stopper
}

func AddEvent(obj interface{}) {
	event := obj.(*corev1.Event)

	// When the event namespace isn't among the namespaces to watch
	if !IsInNamespacesToWatch(event.Namespace, namespacestoWatch) {
		return
	}

	switch cfg.GetString("OUTPUT") {
	case "console":
		fmt.Println(FormatConsole(event))
	case "json":
		b, err := FormatJSON(event)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println(string(b))
	default:
		fmt.Println(FormatConsole(event))
	}
}

func UpdateEvent(objOld interface{}, objNew interface{}) {}

func DeleteEvent(obj interface{}) {}

func IsInNamespacesToWatch(namespace string, namespacesToWatch []string) bool {
	for _, v := range namespacesToWatch {
		if namespace == v {
			return true
		}
	}
	return false
}
