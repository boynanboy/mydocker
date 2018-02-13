package main

import (
	"./container"
	"./cgroups/subsystems"
	"./cgroups"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
)


func Run(tty bool, comArray []string, res *subsystems.ResourceConfig,
         volume string, containerName string, imageName string) {
	if containerName == "" {
		containerName = "noname" + container.RandStringBytes(6)
	}
	parent, writePipe := container.NewParentProcess(tty, volume,
                                                    containerName, imageName)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//record container info
    log.Infof("parent process id: %d", parent.Process.Pid)
	containerName, err := container.RecordContainerInfo(parent.Process.Pid, comArray,
                                              containerName, volume)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
    // todo for test purpose rm the -ti container when its stop
    if tty {
	    parent.Wait()
		//deleteContainerInfo(containerName)
        container.DeleteWorkSpace(containerName, volume)
    }
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}



//func deleteContainerInfo(containerId string) {
	//dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerId)
	//if err := os.RemoveAll(dirURL); err != nil {
		//log.Errorf("Remove dir %s error %v", dirURL, err)
	//}
//}

