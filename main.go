package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/disaster37/migrate-nfs/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var version = "develop"
var commit = ""

func run(args []string) error {

	// Logger setting
	log.SetOutput(os.Stdout)

	// Get home directory
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("Can't get home directory: %s", err.Error())
		homePath = "/root"
	}

	// CLI settings
	app := cli.NewApp()
	app.Usage = "Migrate NFS"
	app.Version = fmt.Sprintf("%s-%s", version, commit)
	app.Flags = []cli.Flag{
		altsrc.NewPathFlag(&cli.PathFlag{
			Name:    "kubeconfig",
			Usage:   "The kube config file",
			EnvVars: []string{"KUBECONFIG"},
			Value:   fmt.Sprintf("%s/.kube/config", homePath),
		}),
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Display debug output",
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Display debug output",
		},
		altsrc.NewInt64Flag(&cli.Int64Flag{
			Name:  "timeout",
			Usage: "The timeout in second",
			Value: 0,
		}),
		&cli.BoolFlag{
			Name:  "no-color",
			Usage: "No print color",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "fix-nfs-version",
			Usage:  "Force NFS version 3 on all NFS PV",
			Action: cmd.FixNFSVersion,
		},
		{
			Name:   "fix-existing-nfsv3",
			Usage:  "Add some option on existing nfsv3",
			Action: cmd.FixNFSV3,
		},
		{
			Name:  "disable-workloads",
			Usage: "Set 0 replica on all deployment and statefullset on provided namespace",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:     "exclude-namespaces",
					Usage:    "The exclude namespaces",
					Required: false,
				},
				&cli.StringSliceFlag{
					Name:     "include-namespaces",
					Usage:    "The include namespaces",
					Required: false,
				},
			},
			Action: cmd.DisableWorkloads,
		},
		{
			Name:  "enable-workloads",
			Usage: "Set 1 replica on all deployment and statefullset on provided namespace",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:     "exclude-namespaces",
					Usage:    "The exclude namespaces",
					Required: false,
				},
				&cli.StringSliceFlag{
					Name:     "include-namespaces",
					Usage:    "The include namespaces",
					Required: false,
				},
			},
			Action: cmd.EnableWorkloads,
		},
		{
			Name:  "enable-workloads",
			Usage: "Set 1 replica on all deployment and statefullset on provided namespace",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:     "exclude-namespaces",
					Usage:    "The exclude namespaces",
					Required: false,
				},
				&cli.StringSliceFlag{
					Name:     "include-namespaces",
					Usage:    "The include namespaces",
					Required: false,
				},
			},
			Action: cmd.EnableWorkloads,
		},
		{
			Name:   "get-nodes-to-reboot",
			Usage:  "Get the list of nodes to reboot",
			Action: cmd.GetNodesToReboot,
		},
	}

	app.Before = func(c *cli.Context) error {

		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		if !c.Bool("no-color") {
			formatter := new(prefixed.TextFormatter)
			formatter.FullTimestamp = true
			formatter.ForceFormatting = true
			log.SetFormatter(formatter)
		}
		return nil
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err = app.Run(args)
	return err
}

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
