package controller

import (
	"fmt"
	"github.com/practice/kube-event/pkg/components/elasticsearchCollector"
	"github.com/practice/kube-event/pkg/components/logCollector"
	"github.com/practice/kube-event/pkg/components/prometheusCollector"
	"github.com/practice/kube-event/pkg/components/sender"
	"github.com/practice/kube-event/pkg/config"
	"github.com/practice/kube-event/pkg/model"
	"k8s.io/klog/v2"
	"time"

	v1 "k8s.io/api/events/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// EventWorker 事件通知器
type EventWorker struct {
	// config 配置文件
	config *config.Config
	// Events 存储event的队列，使用chan实现
	Events chan *model.Event
	stopCh <-chan struct{}
	// clientSet k8s客户端
	clientSet kubernetes.Interface
	// factory informer factory对象
	factory informers.SharedInformerFactory
	// prometheusCollect prometheus收集器
	prometheusCollect *prometheusCollector.PrometheusCollector
	// logCollect 日志收集器
	logCollect *logCollector.StructLogger
	// elasticSearchCollect elastic收集器
	elasticSearchCollect *elasticsearchCollector.ElasticSearchCollector
	// sendCollect 信息发送器
	sendCollect *sender.Sender
}

func NewWorker(stopCh <-chan struct{}, cfg *config.Config) *EventWorker {
	w := &EventWorker{
		Events: make(chan *model.Event, 1000),
		stopCh: stopCh,
		config: cfg,
	}
	w.setMode(cfg)
	return w
}

func (w *EventWorker) setMode(cfg *config.Config) {
	if cfg.Mode.Log {
		w.logCollect = logCollector.NewStructLogger(cfg.LogFilePath)
	}
	if cfg.Mode.Prometheus {
		w.prometheusCollect = prometheusCollector.MetricsCollector
	}
	if cfg.Mode.Message {
		w.sendCollect = sender.NewSender(cfg)
	}
	if cfg.Mode.ElasticSearch {
		w.elasticSearchCollect = elasticsearchCollector.NewElasticSearchCollector(cfg)
	}
}

// initClient 初始化k8s client
func (w *EventWorker) initClient() error {
	cfg, err := clientcmd.BuildConfigFromFlags("", w.config.KubeConfig)
	if err != nil {
		return err
	}
	//cfg.Insecure = true
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	w.clientSet = clientSet
	return nil
}

// Prepare event worker前置准备
func (w *EventWorker) Prepare() error {
	// 初始化k8s client
	err := w.initClient()
	if err != nil {
		return err
	}
	// 初始化 informer
	kubeInformerFactory := informers.NewSharedInformerFactory(w.clientSet, time.Second*30)
	eventInformer := kubeInformerFactory.Events().V1().Events().Informer()
	eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: w.eventAddHandle,
	})
	w.factory = kubeInformerFactory

	return nil
}

// Run 执行event worker
func (w *EventWorker) Run() {
	// 执行factory
	w.factory.Start(w.stopCh)
	w.factory.WaitForCacheSync(w.stopCh)
}

// eventAddHandle event资源对象的回调，主要把监听到的资源对象放入chan中
func (w *EventWorker) eventAddHandle(obj interface{}) {
	// watch add-event, and update events ---> Worker events
	event := obj.(*v1.Event)
	klog.Infof("event: %s, type: %s, reason: %s", event.Name, event.Type, event.Reason)
	var eventTmp model.Event
	eventTmp.Type = event.Type
	eventTmp.Kind = event.Regarding.Kind
	eventTmp.Name = event.Regarding.Name
	eventTmp.Message = event.Note
	eventTmp.Host = event.DeprecatedSource.Host
	eventTmp.Namespace = event.Namespace
	eventTmp.Count = event.DeprecatedCount
	eventTmp.Reason = event.Reason
	eventTmp.Source = event.DeprecatedSource.Component
	eventTmp.Timestamp = event.DeprecatedLastTimestamp.Time

	// TODO: 这里可以做过滤操作
	w.Events <- &eventTmp

}

// Do 执行监听到的event事件
func (w *EventWorker) Do() {
	for {
		select {
		case e := <-w.Events:
			// 从chan 获取event对象，并执行操作

			// 代码写法很丑，可优化，不过能先使用
			if w.config.Mode.Log {
				w.logCollect.EventLog(e)
			}
			if w.config.Mode.Prometheus {
				w.prometheusCollect.Collecting(e)
			}
			if w.config.Mode.Message {
				w.sendCollect.Send(e)
			}
			if w.config.Mode.ElasticSearch {
				// FIXME: 使用goroutine是因为此接口返回很慢
				go func() {
					err := w.elasticSearchCollect.Collecting(e)
					if err != nil {
						fmt.Println("elasticSearchCollect error: ", err)
					}
				}()
			}

		case <-w.stopCh:
			return
		}
	}
}