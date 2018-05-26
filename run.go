package main

import (
	"./container"
	"./cgroups/subsystems"
	"./cgroups"
    "./network"
	log "github.com/Sirupsen/logrus"
	"os"
)


func Run(tty bool, command string, res *subsystems.ResourceConfig,
         volume string, containerName string, imageName string,
         envSlice []string, nw string, portMapping []string) {
	if containerName == "" {
		containerName = "noname" + container.RandStringBytes(6)
	}
	parent, writePipe := container.NewParentProcess(tty, volume, containerName,
                                                    imageName, envSlice)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//record container info
    log.Infof("parent process id: %d", parent.Process.Pid)
	containerInfo, err := container.RecordContainerInfo(parent.Process.Pid, command,
                                                        containerName, volume, portMapping)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

    if nw != "" {
        // config container network
        network.Init()
        if err := network.Connect(nw, containerInfo); err != nil {
            // this is where is wrong
            log.Errorf("Error Connect Network: %v", err)
            return
        }
    }

	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager(containerName)
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)


	sendInitCommand(command, writePipe)
    // todo for test purpose rm the -ti container when its stop
    if tty {
	    parent.Wait()
        container.DeleteWorkSpace(containerName, volume)
    }
}

func sendInitCommand(command string, writePipe *os.File) {
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

