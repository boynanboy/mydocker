package main

import (
	"fmt"
    "os"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"./cgroups/subsystems"
	"./container"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit mydocker run [-ti|-d] [command]`,
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
        imageName := cmdArray[0]
        cmdArray = cmdArray[1:]
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
        log.Infof("tty enabled %v", tty)
		volume := context.String("v")
		containerName := context.String("name")
		Run(tty, cmdArray, resConf, volume, containerName, imageName)
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


var logCommand = cli.Command{
	Name: "logs",
	Usage: "print logs of a container",
	Action: func(context *cli.Context) error {
        if len(context.Args()) < 1 {
            return fmt.Errorf("Please input your container name")
        }
        containerName := context.Args().Get(0)
        logContainer(containerName)
        return nil
	},
}

var execCommand = cli.Command{
	Name: "exec",
	Usage: "exec a command into container",
	Action: func(context *cli.Context) error {
		//This is for callback
		if os.Getenv(ENV_EXEC_PID) != "" {
			log.Infof("pid callback pid %s", os.Getgid())
			return nil
		}

		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container name or command")
		}
		containerName := context.Args().Get(0)
		var commandArray []string
		for _, arg := range context.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		ExecContainer(containerName, commandArray)
		return nil
 	},
}

var stopCommand = cli.Command{
	Name: "stop",
	Usage: "stop a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		stopContainer(containerName)
		return nil
	},
}

var removeCommand = cli.Command{
	Name: "rm",
	Usage: "remove unused containers",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
        // TODO this remove container may not work properly when container have
        // volume have not tested it
	    container.RemoveContainer(containerName)
		return nil
	},
}

// This is just a temporary command for clean up all containers
// (make sure all containers stop firs) since "mydocker ps -q" does not exsits
var cleanCommand = cli.Command{
	Name: "clean",
	Usage: "clean all the containers (for easy development) have to stop",
	Action: func(context *cli.Context) error {
        mergedDir := "./merged/"
        writeLayerDir := "./container_layer/"
        readLayerDir := "./base/"
        indexDir := "./index/"
        logDir := "./logs/"
        infoDir := "./info/"

        // remove all the storage of containers
        if err := os.RemoveAll(mergedDir); err != nil {
            log.Errorf("Remove dir %s error %v", mergedDir, err)
        }
        if err := os.MkdirAll(mergedDir, 0622); err != nil {
            log.Errorf("Mkdir error %s error %v", mergedDir, err)
        }
        if err := os.RemoveAll(writeLayerDir); err != nil {
            log.Errorf("Remove dir %s error %v", writeLayerDir, err)
        }
        if err := os.MkdirAll(writeLayerDir, 0622); err != nil {
            log.Errorf("Mkdir error %s error %v", writeLayerDir, err)
        }
        if err := os.RemoveAll(indexDir); err != nil {
            log.Errorf("Remove dir %s error %v", indexDir, err)
        }
        if err := os.MkdirAll(indexDir, 0622); err != nil {
            log.Errorf("Mkdir error %s error %v", indexDir, err)
        }
        // remove container info and log if there is any
        if err := os.RemoveAll(logDir); err != nil {
            log.Errorf("Remove dir %s error %v", logDir, err)
        }
        if err := os.MkdirAll(logDir, 0622); err != nil {
            log.Errorf("Mkdir error %s error %v", logDir, err)
        }
        if err := os.RemoveAll(infoDir); err != nil {
            log.Errorf("Remove dir %s error %v", infoDir, err)
        }
        if err := os.MkdirAll(infoDir, 0622); err != nil {
            log.Errorf("Mkdir error %s error %v", infoDir, err)
        }
        if err := os.RemoveAll(readLayerDir); err != nil {
            log.Errorf("Remove dir %s error %v", infoDir, err)
        }
        if err := os.MkdirAll(readLayerDir, 0622); err != nil {
            log.Errorf("Mkdir error %s error %v", infoDir, err)
        }

        return nil
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container into image",
 	Action: func(context *cli.Context) error {
		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container name and image name")
		}
		containerName := context.Args().Get(0)
		imageName := context.Args().Get(1)
		commitContainer(containerName, imageName)
		return nil
    },
}
