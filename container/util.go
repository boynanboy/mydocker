package container

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
    "strconv"
    "time"
    "math/rand"
	"encoding/json"
)


func RecordContainerInfo(containerPID int, command string,
                         containerName string, volume string,
                         portMapping []string) (*ContainerInfo, error) {
	id := RandStringBytes(12)
	createTime := time.Now().Format("2006-01-02 15:04:05")
	//createTime := time.Now()
	if containerName == "" {
		containerName = id
	}
	containerInfo := &ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      RUNNING,
		Name:        containerName,
		Volume:      volume,
		PortMapping: portMapping,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return nil, err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirUrl, err)
		return nil, err
	}
	fileName := dirUrl + "/" + ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return nil, err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
        return nil, err
	}

	return containerInfo, nil
}

func GetContainerInfo(file os.FileInfo) (*ContainerInfo, error) {
	containerName := file.Name()
	configFileDir := fmt.Sprintf(DefaultInfoLocation, containerName)
	configFileDir = configFileDir + ConfigName
	content, err := ioutil.ReadFile(configFileDir)
	if err != nil {
		log.Errorf("Read file %s error %v", configFileDir, err)
		return nil, err
	}
	var containerInfo ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Json unmarshal error %v", err)
		return nil, err
	}

	return &containerInfo, nil
}

func GetContainerPidByName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	configFilePath := dirURL + ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	var containerInfo ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", err
	}
	return containerInfo.Pid, nil
}

func GetContainerInfoByName(containerName string) (*ContainerInfo, error) {
	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	configFilePath := dirURL + ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Errorf("Read file %s error %v", configFilePath, err)
		return nil, err
	}
	var containerInfo ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		log.Errorf("GetContainerInfoByName unmarshal error %v", err)
		return nil, err
	}
	return &containerInfo, nil
}

func RemoveContainer(containerName string) {
	containerInfo, err := GetContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	if containerInfo.Status != STOP {
		log.Errorf("Couldn't remove running container")
		return
	}
    DeleteWorkSpace(containerName, "")
}

func GetEnvsByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Read file %s error %v", path, err)
		return nil
	}
	//env split by \u0000 ?
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs
}

func RandStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
