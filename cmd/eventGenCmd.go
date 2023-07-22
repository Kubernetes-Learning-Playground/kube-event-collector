package cmd

import (
	"github.com/practice/kube-event/pkg/generator"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func kubeEventGeneratorCmd() *cobra.Command {
	eg := generator.NewGenerator()
	cmd := &cobra.Command{
		Use:   "kube-event-generator",
		Short: "run kube-event-generator",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			if err := eg.Run(); err != nil {
				klog.Error(err)
			}
		},
	}

	eg.AddFlags(cmd)
	return cmd
}
