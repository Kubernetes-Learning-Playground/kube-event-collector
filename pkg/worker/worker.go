package worker

import (
	"github.com/practice/kube-event/pkg/config"
	"github.com/practice/kube-event/pkg/server"
	"github.com/practice/kube-event/pkg/signal"
	"github.com/practice/kube-event/pkg/worker/controller"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// EventWorker 处理事件对象
type EventWorker struct {
	config     string
	serverPort string
}

func NewEventWorker() *EventWorker {
	return &EventWorker{}
}

// AddFlags adds flags to command.
func (e *EventWorker) AddFlags(cmd *cobra.Command) {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(cmd.PersistentFlags())

	flags := cmd.Flags()
	flags.StringVar(&e.config, "config1", "./config1.yaml", "config1")
	flags.StringVar(&e.serverPort, "port", "8080", "server port ")

}

func (e *EventWorker) Execute() error {
	// 配置文件
	c, err := config.LoadConfig(e.config)
	if err != nil {
		return err
	}

	stopC := signal.SetupSignalHandler()
	worker := controller.NewWorker(stopC, c)

	// 1. worker 初始化操作
	err = worker.Prepare()
	if err != nil {
		return err
	}

	// 2. 启动worker中的informer
	worker.Run()

	// 3. 不断轮巡，进行collector操作
	go worker.Do()

	// 启动http server，主要给prometheus调用接口
	server.HttpServer(e.serverPort)
	return nil
}