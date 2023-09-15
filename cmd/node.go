package cmd

import (
	"strings"

	"github.com/disaster37/migrate-nfs/pkg/k8s"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func GetNodesToReboot(c *cli.Context) (err error) {

	client, err := getClient(c)
	if err != nil {
		return err
	}

	ctx, ctxCancel := getContext(c)
	if ctxCancel != nil {
		defer ctxCancel()
	}

	nodes, err := k8s.GetNodesWhereNfsIsMounted(ctx, client)
	if err != nil {
		return err
	}

	log.Infof("You need to reboot the nodes: %s", strings.Join(nodes, ", "))

	return nil
}
