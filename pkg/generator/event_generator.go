package generator

import (
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/reference"
	"k8s.io/klog/v2"
	"time"
)

// EventGenerator generates event.
type EventGenerator struct {
	kubeConfig string

	kind      string
	name      string
	namespace string

	eventType    string
	eventAction  string
	eventReason  string
	eventMessage string

	restClientGetter resource.RESTClientGetter
}

const (
	defaultNamespace = "kube-system"
	defaultEventType = v1.EventTypeWarning
)

// NewGenerator creates a generator.
func NewGenerator() *EventGenerator {
	return &EventGenerator{}
}

// AddFlags adds flags to command.
func (g *EventGenerator) AddFlags(cmd *cobra.Command) {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(cmd.PersistentFlags())
	g.restClientGetter = configFlags

	flags := cmd.Flags()
	flags.StringVar(&g.kubeConfig, "kubeconfig", "/Users/zhenyu.jiang/.kube/config", "kube config file")
	flags.StringVar(&g.kind, "kind", "pods", "Resource kind to get.")
	flags.StringVar(&g.name, "name", "kube-scheduler-minikube", "Resource name to get.")
	flags.StringVar(&g.namespace, "namespace", defaultNamespace, "Resource namespace to get.")
	flags.StringVar(&g.eventType, "type", defaultEventType, "Event type.")
	flags.StringVar(&g.eventAction, "action", "ttt", "Event action.")
	flags.StringVar(&g.eventReason, "reason", "Testing-Reason", "Event reason.")
	flags.StringVar(&g.eventMessage, "message", "Testing-Message", "Event message.")
}

// Run generates event.
func (g *EventGenerator) Run() error {
	r := resource.NewBuilder(g.restClientGetter).
		Unstructured().
		NamespaceParam(g.namespace).
		ResourceTypeOrNameArgs(true, g.kind, g.name).
		Do()
	if err := r.Err(); err != nil {
		return err
	}

	infos, err := r.Infos()
	if err != nil {
		return err
	}

	ref, err := reference.GetReference(scheme.Scheme, infos[0].Object)
	if err != nil {
		return err
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", g.kubeConfig)
	if err != nil {
		return err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	if len(g.eventAction) == 0 {
		g.eventAction = g.eventReason
	}

	now := time.Now()
	event, err := client.CoreV1().Events("").CreateWithEventNamespace(&v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v.%x", g.name, now.UnixNano()),
			Namespace: g.namespace,
		},
		FirstTimestamp:      metav1.NewTime(now),
		LastTimestamp:       metav1.NewTime(now),
		EventTime:           metav1.NewMicroTime(now),
		ReportingController: "k8s-event-generator",
		ReportingInstance:   "k8s-event-generator",
		Action:              g.eventAction,
		InvolvedObject:      *ref,
		Reason:              g.eventReason,
		Type:                g.eventType,
		Message:             g.eventMessage,
	})

	if err == nil {
		klog.Infof("Event generated successfully: %v", event)
	}
	return  nil
}
