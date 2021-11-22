package dto

import "bryson.foundation/kbuildresource/models"

// 少了一些非必要参数，多了一些业务属性配置，比如tuning
type BuildJobDTO struct {
	ClusterName string `json:"clusterName" description:"必选，集群名"`
	NetworkZone string `json:"networkZone" description:"网络区域"`
	Name string `json:"name" description:"等于slavename"`
	ReName bool `json:"reName" description:"接受重命名"`
	Labels []string `json:"labels" description:"标签"`
	Namespace string `json:"namespace" description:"命名空间"`
	Tuning bool `json:"tuning" description:"是否接受资源参数优化"`
	Containers []*models.Container `json:"containers" description:"容器配置"`
	InstanceName string `json:"instance_name"`
}

//type ContainerDTO struct {
//	CMDs []string `json:"cmd" description:"容器启动命令：eg: cm1,cmd2"`
//	Name string `json:"name" description:"容器名"`
//	Image string `json:"image" description:"容器镜像"`
//	RequestCPU string`json:"requestCPU" description:"请求cpu大小"`
//	RequestMem string`json:"requestMem" description:"请求内存大小"`
//	LimitCPU string`json:"requestCPU" description:"最大可用cpu大小"`
//	LimitMem string`json:"requestCPU" description:"最大可用内存大小"`
//}
