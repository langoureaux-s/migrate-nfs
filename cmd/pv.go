package cmd

import (
	"github.com/disaster37/migrate-nfs/pkg/k8s"
	"github.com/urfave/cli/v2"
)

func FixNFSVersion(c *cli.Context) (err error) {

	client, err := getClient(c)
	if err != nil {
		return err
	}

	ctx, ctxCancel := getContext(c)
	if ctxCancel != nil {
		defer ctxCancel()
	}

	return k8s.MigrateNFSPV(ctx, client, c.Bool("dry-run"))
}

func FixNFSV3(c *cli.Context) (err error) {

	client, err := getClient(c)
	if err != nil {
		return err
	}

	ctx, ctxCancel := getContext(c)
	if ctxCancel != nil {
		defer ctxCancel()
	}

	return k8s.FixNFSPV(ctx, client, c.Bool("dry-run"))
}
