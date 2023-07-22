package worker

import (
	"github.com/practice/kube-event/pkg/config"
	"github.com/practice/kube-event/pkg/server"
	"github.com/practice/kube-event/pkg/signal"
	"github.com/practice/kube-event/pkg/worker/controller"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type EventWorker struct {
	config string
}

func NewEventWorker() *EventWorker {
	return &EventWorker{}
}

// AddFlags adds flags to command.
func (e *EventWorker) AddFlags(cmd *cobra.Command) {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(cmd.PersistentFlags())

	flags := cmd.Flags()
	flags.StringVar(&e.config, "config", "./config.yaml", "config ")

}

func (e *EventWorker) Execute() error {
	c, err := config.LoadConfig(e.config)
	if err != nil {
		return err
	}

	stopC := signal.SetupSignalHandler()
	w := controller.NewWorker(stopC, c)
	err = w.Prepare()
	if err != nil {
		return err
	}
	w.Run()
	go w.Do()

	// 启动http server，主要给prometheus调用接口
	server.HttpServer()
	return nil
}
