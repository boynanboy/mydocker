package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"./container"
	"os/exec"
)

func commitContainer(containerName, imageName string) {
	mergedURL := fmt.Sprintf(container.DefaultMergedLocation, containerName)
	imageTar := fmt.Sprintf(container.DefaultImagesLocation, containerName)

    log.Infof("Tar folder merged url: %s, tar the url: %s", mergedURL, imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mergedURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mergedURL, err)
	}
}
