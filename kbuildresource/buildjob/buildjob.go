package buildjob

import (
	"bryson.foundation/kbuildresource/dto"
	"bryson.foundation/kbuildresource/models"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func CreateBuildJob(buildJobDTO *dto.BuildJobDTO) error {
	fmt.Println("Before executing buildJob")
	// 对于DTO的参数校验
	err := VerifyBuildJobDTO(buildJobDTO)
	if err != nil {
		return err
	}
	err = CreatePod(buildJobDTO)
	if err != nil {
		return err
	}
	return nil
}

func VerifyBuildJobDTO(buildJobDTO *dto.BuildJobDTO) error {
	logrus.Info("INFO: verifyBuildJobDTO")
	//time.Sleep(1 * time.Second)
	return nil
}

func CreatePod(buildJobDTO *dto.BuildJobDTO) error {
	logrus.Info("INFO: CreatePod")
	_, err := createPodFromBuildJobDTO(buildJobDTO)
	if err != nil {
		logrus.Info("INFO: finish CreatePod")
		return nil
	}
	return err
}

func createPodFromBuildJobDTO(buildJobDTO *dto.BuildJobDTO) (*models.Pod, error) {
	pod := &models.Pod{
		Name:        buildJobDTO.Name,
		ClusterName: buildJobDTO.ClusterName,
		Labels:      strings.Join(buildJobDTO.Labels, ","),
		Namespace:   buildJobDTO.Namespace,
		Status:      "Pending",
		NodeIP:      "",
		IsDelete:    "0",
		Containers:  buildJobDTO.Containers,
	}
	_, err := models.AddPod(pod)
	if err != nil {
		logrus.Error("ERROR: add pod to mysql failed, error: ", err)
	}
	time.Sleep(1 * time.Second)
	logrus.Info("INFO: finish create pod")
	return pod, nil
}