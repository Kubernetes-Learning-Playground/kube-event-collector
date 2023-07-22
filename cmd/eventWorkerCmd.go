package cmd

import (
	"github.com/practice/kube-event/pkg/worker"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func kubeEventWorkerCmd() *cobra.Command {
	ew := worker.NewEventWorker()
	cmd := &cobra.Command{
		Use:   "kube-event-worker",
		Short: "run kube-event-worker",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			if err := ew.Execute(); err != nil {
				klog.Error(err)
			}
		},
	}
	ew.AddFlags(cmd)

	return cmd
}
