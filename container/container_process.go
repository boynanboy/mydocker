package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
    "fmt"
	"syscall"
	"strings"
)

var (
	RUNNING             string = "running"
	STOP                string = "stopped"
	Exit                string = "exited"
	DefaultInfoLocation string = "./info/%s/"
	DefaultLogLocation string = "./logs/%s.log"
	DefaultMergedLocation string = "./merged/%s/"
	DefaultImagesLocation string = "./images/%s.tar/"
	DefaultReadOnlyLocationLayer string = "./base/%s/"
	DefaultWritableLayerLocation string = "./container_layer/%s/"
	ConfigName          string = "config.json"
)

type ContainerInfo struct {
	Pid         string `json:"pid"` //容器的init进程在宿主机上的 PID
	Id          string `json:"id"`  //容器Id
	Name        string `json:"name"`  //容器名
	Command     string `json:"command"`    //容器内init运行命令
	CreatedTime string `json:"createTime"` //创建时间
	Status      string `json:"status"`     //容器的状态
	Volume      string `json:"volume"`     //容器的mounted volume
}

func NewParentProcess(tty bool, volume , containerName, imageName string) (*exec.Cmd, *os.File) {
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
	} else {
        // makes logs dir if it is not exists
        os.MkdirAll("./logs/", os.ModePerm)
        stdLogFilePath := fmt.Sprintf(DefaultLogLocation, containerName)
        stdLogFile, err := os.Create(stdLogFilePath)
        if err != nil {
            log.Errorf("NewParentProcess failed to create log file %s error %v",
                       stdLogFilePath, err)
            return nil, nil
        }
        cmd.Stdout = stdLogFile
 	}

	cmd.ExtraFiles = []*os.File{readPipe}

    // a index thing that is only needed for overlayfs do not totally
    // understand yet
	NewWorkSpace(imageName, volume, containerName)
	cmd.Dir = fmt.Sprintf(DefaultMergedLocation, containerName)
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

// create a temp read only layer
// this is implementation is a extreme naive one and does not follow how docker
// really works, since this is a proof of concept project to help me to understand
// how docker works
func createReadOnlyLayer(imageName string) error {
	unTarFolderUrl := fmt.Sprintf(DefaultReadOnlyLocationLayer, imageName)
	imageUrl := fmt.Sprintf(DefaultImagesLocation, imageName)
	exist, err := PathExists(unTarFolderUrl)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", unTarFolderUrl, err)
		return err
	}
	if !exist {
		if err := os.MkdirAll(unTarFolderUrl, 0622); err != nil {
			log.Errorf("Mkdir %s error %v", unTarFolderUrl, err)
			return err
		}

		if _, err := exec.Command("tar", "-xvf", imageUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", unTarFolderUrl, err)
			return err
		}
	}
	return nil
}

func createContainerLayer(mergedURL string, imageName string, indexURL string, writeLayerURL string) {
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
	baseURL := fmt.Sprintf(DefaultReadOnlyLocationLayer, imageName)

    dirs := "lowerdir=" + baseURL + ",upperdir=" + writeLayerURL + ",workdir=" + indexURL
    log.Infof("overlayfs union parameters: %s", dirs)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs,  mergedURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}

}

// create a temp read only layer
func CreateReadOnlyLayer(imageName string) error {
	unTarFolderUrl := fmt.Sprintf(DefaultReadOnlyLocationLayer, imageName)
	imageUrl := fmt.Sprintf(DefaultImagesLocation, imageName)
	exist, err := PathExists(unTarFolderUrl)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", unTarFolderUrl, err)
		return err
	}
	if !exist {
		if err := os.MkdirAll(unTarFolderUrl, 0622); err != nil {
			log.Errorf("Mkdir %s error %v", unTarFolderUrl, err)
			return err
		}

		if _, err := exec.Command("tar", "-xvf", imageUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", unTarFolderUrl, err)
			return err
		}
	}
	return nil
}

func NewWorkSpace(imageName string, volume string, containerName string) {
    mergedURL := "./merged/" + containerName
    indexURL := "./index/" + containerName
    writeLayerURL := "./container_layer/" + containerName
    createReadOnlyLayer(imageName)
    createContainerLayer(mergedURL, imageName, indexURL, writeLayerURL)
    if(volume != ""){
        volumeURLs := volumeUrlExtract(volume)
        length := len(volumeURLs)
        if(length == 2 && volumeURLs[0] != "" && volumeURLs[1] !=""){
            log.Infof("%q",volumeURLs)
            MountVolume(mergedURL, volumeURLs)
        } else {
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

func DeleteWorkSpace(containerName string, volume string) {
    mergedURL := "./merged/" + containerName
    writeLayerURL := "./container_layer/" + containerName
    indexURL := "./index/" + containerName
    logURL := fmt.Sprintf(DefaultLogLocation, containerName)
    infoURL := fmt.Sprintf(DefaultInfoLocation, containerName)
    if (volume == "") {
	    containerInfo, _ := GetContainerInfoByName(containerName)
        volume = containerInfo.Volume
    }
    if (volume != "") {
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
    // remove container info and log if there is any
	if err := os.RemoveAll(logURL); err != nil {
		log.Errorf("Remove dir %s error %v", logURL, err)
	}
	if err := os.RemoveAll(infoURL); err != nil {
		log.Errorf("Remove dir %s error %v", infoURL, err)
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
