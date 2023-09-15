package cmd

import (
	"github.com/disaster37/migrate-nfs/pkg/k8s"
	"github.com/urfave/cli/v2"
)

func DisableWorkloads(c *cli.Context) (err error) {

	client, err := getClient(c)
	if err != nil {
		return err
	}

	ctx, ctxCancel := getContext(c)
	if ctxCancel != nil {
		defer ctxCancel()
	}

	return k8s.DisableWorkLoadNFS(ctx, client, c.StringSlice("exclude-namespaces"), c.StringSlice("include-namespaces"), c.Bool("dry-run"))
}

func EnableWorkloads(c *cli.Context) (err error) {
	client, err := getClient(c)
	if err != nil {
		return err
	}

	ctx, ctxCancel := getContext(c)
	if ctxCancel != nil {
		defer ctxCancel()
	}

	return k8s.EnableWorkLoadNFS(ctx, client, c.StringSlice("exclude-namespaces"), c.StringSlice("include-namespaces"), c.Bool("dry-run"))
}
