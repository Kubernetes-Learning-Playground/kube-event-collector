package prometheusCollector

import (
	"github.com/practice/kube-event/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsCollector prometheus收集器
var MetricsCollector *PrometheusCollector

func init() {
	MetricsCollector = NewPrometheusCollector()
}

// PrometheusCollector prometheus收集器
type PrometheusCollector struct {
	EventCounterVec *prometheus.CounterVec
	//OrderCounterVec   *prometheus.CounterVec
	//RequestCounterVec *prometheus.CounterVec
	//GaugeVec          *prometheus.GaugeVec
}

// NewPrometheusCollector prometheus collector
func NewPrometheusCollector() *PrometheusCollector {
	return &PrometheusCollector{
		EventCounterVec: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "k8s_event_count",
			Help: "The total number of processed events",
		}, []string{"type", "kind", "name", "namespace", "message", "reason", "source", "host"}), // 这里必须写全部的

	}
}

func (pm *PrometheusCollector) Collecting(event *model.Event) {
	pm.EventCounterVec.With(
		prometheus.Labels{
			"type":      event.Type,
			"kind":      event.Kind,
			"name":      event.Name,
			"namespace": event.Namespace,
			"message":   event.Message,
			"reason":    event.Reason,
			"source":    event.Source,
			"host":      event.Host,
			//"timestamp": event.Timestamp.String(),  // 不要记录时间，因为需要记录同一个pod的次数
		},
	).Inc()
}
