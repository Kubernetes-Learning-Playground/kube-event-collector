package elasticsearchCollector

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"github.com/practice/kube-event/pkg/config"
	"github.com/practice/kube-event/pkg/model"
	"k8s.io/klog/v2"
	"log"
)

// ElasticSearchCollector es收集器
type ElasticSearchCollector struct {
	EsClient *elastic.Client
}

func NewElasticSearchCollector(config *config.Config) *ElasticSearchCollector {
	ec, err := initElasticSearch(config.ElasticSearchEndpoint)
	if err != nil {
		panic(err)
	}
	return &ElasticSearchCollector{
		EsClient: ec,
	}
}

func initElasticSearch(endpoint string) (*elastic.Client, error) {

	es, err := elastic.NewClient(
		elastic.SetURL(endpoint),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	_, _, err = es.Ping(endpoint).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return es, nil
}

func (ec *ElasticSearchCollector) Collecting(event *model.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
		return err
	}

	d, err := ec.EsClient.Index().Index("events").
		Id("").BodyString(string(data)).Do(context.Background())

	klog.Infof("插入数据完成  %v\n", d.Id)

	return nil
}
