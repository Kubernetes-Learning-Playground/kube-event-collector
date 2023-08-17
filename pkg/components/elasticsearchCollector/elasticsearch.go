package elasticsearchCollector

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/practice/kube-event/pkg/config"
	"github.com/practice/kube-event/pkg/model"
	"k8s.io/klog/v2"
	"log"
)

// ElasticSearchCollector es收集器
type ElasticSearchCollector struct {
	EsClient *elasticsearch.Client
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

func initElasticSearch(endpoint string) (*elasticsearch.Client, error) {
	var err error
	var es *elasticsearch.Client
	cfg := elasticsearch.Config{
		Addresses: []string{
			endpoint,
		},
		//Username: eventFlag.UserName,
		//Password: eventFlag.Password,
		//		CACert: cert,
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	//_, err = es.Info()
	//if err != nil {
	//	panic(err.Error())
	//}
	return es, nil
}

func (ec *ElasticSearchCollector) Collecting(event *model.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
		return err
	}
	req := esapi.IndexRequest{
		Index:      "event",
		DocumentID: event.Name,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), ec.EsClient)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()
	klog.Infof("插入数据完成  %s\n", event.Name)
	return nil
}
