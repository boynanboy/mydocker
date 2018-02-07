package main

import (
	"./container"
	"./cgroups/subsystems"
	"./cgroups"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
)


func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//record container info
    log.Infof("parent process id: %d", parent.Process.Pid)

	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
    // stop container will not cause temp container workspace got removed
    // that is done by remove container
    if tty {
	    parent.Wait()
    }
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
