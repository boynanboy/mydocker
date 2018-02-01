package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"./container"
	"os"
	"io/ioutil"
)

func logContainer(containerName string) {
	logFileLocation := fmt.Sprintf(container.DefaultLogLocation, containerName)
	file, err := os.Open(logFileLocation)
	defer file.Close()
	if err != nil {
		log.Errorf("Log container open file %s error %v", logFileLocation, err)
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Log container read file %s error %v", logFileLocation, err)
		return
	}
	fmt.Fprint(os.Stdout, string(content))
}
