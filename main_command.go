package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"./cgroups/subsystems"
	"./container"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run [-ti|-d] [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name: "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name: "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},

	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
        resConf := &subsystems.ResourceConfig{
            MemoryLimit: context.String("m"),
            CpuSet: context.String("cpuset"),
            CpuShare:context.String("cpushare"),
        }

		tty := context.Bool("ti")
        detach := context.Bool("d")
        if tty && detach {
            return fmt.Errorf("ti and d paramter can not both provided")
        }
		volume := context.String("v")
		containerName := context.String("name")
		Run(tty, cmdArray, resConf, volume, containerName)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}

var exportCommand = cli.Command{
	Name:  "export",
	Usage: "export current running container into tar",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		imageName := context.Args().Get(0)
		exportContainer(imageName)
		return nil
	},
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all the containers",
	Action: func(context *cli.Context) error {
		ListContainers()
 		return nil
 	},
 }
