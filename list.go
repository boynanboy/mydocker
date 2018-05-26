package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"./container"
	"io/ioutil"
    "os"
	"text/tabwriter"
)

func ListContainers() {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, "")
	dirURL = dirURL[:len(dirURL)-1]
	files, err := ioutil.ReadDir(dirURL)
	if err != nil {
		log.Errorf("Read dir %s error %v", dirURL, err)
		return
	}

	var containers []*container.ContainerInfo
	for _, file := range files {
		if file.Name() == "network" {
			continue
		}
		tmpContainer, err := container.GetContainerInfo(file)
		if err != nil {
			log.Errorf("Get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}
