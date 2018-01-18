package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
	"strings"
)

func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
    imageURL := "./image"
    // a index thing that is only needed for overlayfs do not totally 
    // understand yet
	NewWorkSpace(imageURL, volume)
	cmd.Dir = "./merged"
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func createContainerLayer(mergedURL string, imageURL string, indexURL string, writeLayerURL string) {
    // for easy coding did not check whether certain folders exists before
    // ideally should do it
	if err := os.Mkdir(writeLayerURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeLayerURL, err)
	}
	if err := os.Mkdir(mergedURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mergedURL, err)
	}
	if err := os.Mkdir(indexURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", indexURL, err)
	}

    dirs := "lowerdir=" + imageURL + ",upperdir=" + writeLayerURL + ",workdir=" + indexURL
    log.Infof("overlayfs union parameters: %s", dirs)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs,  mergedURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}

}

func NewWorkSpace(imageURL string, volume string) {
    mergedURL := "./merged"
    indexURL := "./index"
    writeLayerURL := "./container_layer"
    createContainerLayer(mergedURL, imageURL, indexURL, writeLayerURL)
    if(volume != ""){
        volumeURLs := volumeUrlExtract(volume)
        length := len(volumeURLs)
        if(length == 2 && volumeURLs[0] != "" && volumeURLs[1] !=""){
            log.Infof("%q",volumeURLs)
            MountVolume(mergedURL, volumeURLs)
        }else{
            log.Infof("Volume parameter input is not correct.")
        }
    }
}

func MountVolume(mergedURL string, volumeURLs []string)  {
    parentUrl := volumeURLs[0]
    if err := os.Mkdir(parentUrl, 0777); err != nil {
        // ideally if the system directory is already exists there is no need
        // try to create that, I just do not want to polish this thing here
        log.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
    }
    containerRelativeVolumeUrl := volumeURLs[1]
    containerVolumeURL := mergedURL + containerRelativeVolumeUrl
    if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
        log.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
    }
    cmd := exec.Command("mount", "--bind", parentUrl, containerVolumeURL)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Errorf("Mount volume failed. %v", err)
    }

}


func DeleteWorkSpace(volume string) {
    mergedURL := "./merged"
    writeLayerURL := "./container_layer"
    indexURL := "./index"
    if(volume != ""){
        volumeURLs := volumeUrlExtract(volume)
        length := len(volumeURLs)
        if(length == 2 && volumeURLs[0] != "" && volumeURLs[1] !=""){
            containerRelativeVolumeUrl := volumeURLs[1]
            containerVolumeURL := mergedURL + containerRelativeVolumeUrl
            //DeleteMountPointWithVolume(rootURL, mntURL, volumeURLs)
            cmd := exec.Command("umount", containerVolumeURL)
            cmd.Stdout=os.Stdout
            cmd.Stderr=os.Stderr
            if err := cmd.Run(); err != nil {
                log.Errorf("Umount volume failed. %v",err)
            }
        }
    }
	cmd := exec.Command("umount", mergedURL)

	cmd.Stdout=os.Stdout
	cmd.Stderr=os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v",err)
	}
    // remove merged, index and container write layer
	if err := os.RemoveAll(mergedURL); err != nil {
		log.Errorf("Remove dir %s error %v", mergedURL, err)
	}
	if err := os.RemoveAll(writeLayerURL); err != nil {
		log.Errorf("Remove dir %s error %v", writeLayerURL, err)
	}
	if err := os.RemoveAll(indexURL); err != nil {
		log.Errorf("Remove dir %s error %v", indexURL, err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func volumeUrlExtract(volume string)([]string){
	var volumeURLs []string
	volumeURLs =  strings.Split(volume, ":")
	return volumeURLs
}
