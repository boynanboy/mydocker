package main


import (
	log "github.com/Sirupsen/logrus"
	"fmt"
	"os/exec"
)

func exportContainer(imageName string){
	mergedURL := "./merged/" + imageName
	imageTar := "./" + imageName + ".tar"
	fmt.Printf("exported %s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mergedURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mergedURL, err)
	}
}
