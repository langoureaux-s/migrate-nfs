package cmd

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getClient(c *cli.Context) (client kubernetes.Interface, err error) {

	kubeConfigPath := c.Path("kubeconfig")

	log.Debugf("Use kubeConfig %s", kubeConfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	clientTmp, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientTmp, nil
}

func getContext(c *cli.Context) (ctx context.Context, cancelFunc context.CancelFunc) {
	if c.Int64("timeout") == 0 {
		return c.Context, nil
	}

	return context.WithTimeout(c.Context, time.Duration(c.Int64("timeout"))*time.Second)
}
