package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
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
    //cmd.Dir = "./root/busybox"
	cmd.ExtraFiles = []*os.File{readPipe}
	//mntURL := "./root/mnt/"
	//rootURL := "./root/"
    // at this stage just use easy busybox as the image
    // a more comprehensive image management can be added
    // at this stage image url and write layer and mount point are all
    // static
    imageURL := "./image"
    // a index thing that is only needed for overlayfs do not totally 
    // understand yet
	NewWorkSpace(imageURL)
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

//Create a AUFS filesystem as container root workspace
//func NewWorkSpace(imageURL string) {
    // no need to create read only layer which is already created
	//CreateReadOnlyLayer(rootURL)

    // no needed as independent function
	//CreateWriteLayer(writeLayerURL)
	//CreateMountPoint(imageURL, writeLayerURL, mntURL, indexURL)
//}

//func CreateReadOnlyLayer(rootURL string) {
	//busyboxURL := rootURL + "busybox/"
	//busyboxTarURL := rootURL + "busybox.tar"
	//exist, err := PathExists(busyboxURL)
	//if err != nil {
		//log.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
	//}
	//if exist == false {
		//if err := os.Mkdir(busyboxURL, 0777); err != nil {
			//log.Errorf("Mkdir dir %s error. %v", busyboxURL, err)
		//}
		//if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			//log.Errorf("Untar dir %s error %v", busyboxURL, err)
		//}
	//}
//}

//func CreateWriteLayer(writeLayerURL) {
	//if err := os.Mkdir(writeLayerURL, 0777); err != nil {
		//log.Errorf("Mkdir dir %s error. %v", writeURL, err)
	//}
//}

func NewWorkSpace(imageURL string) {
    mergedURL := "./merged"
    indexURL := "./index"
    writeLayerURL := "./container_layer"
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
    // sudo mount -t overlay overlay -o lowerdir=./image_layer1:./image_layer2,upperdir=./container_layer,workdir=./work ./merged

	//dirs := "dirs=" + rootURL + "writeLayer:" + rootURL + "busybox"
	//cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)

    dirs := "lowerdir=" + imageURL + ",upperdir=" + writeLayerURL + ",workdir=" + indexURL
    log.Infof("overlayfs union parameters: %s", dirs)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs,  mergedURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

// the overlayfs created content for new created container
func DeleteWorkSpace() {
    mergedURL := "./merged"
    writeLayerURL := "./container_layer"
    indexURL := "./index"

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

//func DeleteMountPoint(rootURL string, mntURL string){
    //mergedURL := "./mnt"
	//cmd := exec.Command("umount", mergedURL)
	//cmd.Stdout=os.Stdout
	//cmd.Stderr=os.Stderr
	//if err := cmd.Run(); err != nil {
		//log.Errorf("%v",err)
	//}
	//if err := os.RemoveAll(mergedURL); err != nil {
		//log.Errorf("Remove dir %s error %v", mergedURL, err)
	//}
//}

//func DeleteWriteLayer(rootURL string) {
	//writeURL := rootURL + "writeLayer/"
	//if err := os.RemoveAll(writeURL); err != nil {
		//log.Errorf("Remove dir %s error %v", writeURL, err)
	//}
//}

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
