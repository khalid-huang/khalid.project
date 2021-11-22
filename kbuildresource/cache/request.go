package cache

import (
	"bryson.foundation/kbuildresource/common"
	"bryson.foundation/kbuildresource/models"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/prometheus/common/log"
)

var (
	instanceNameListKey = GenInstanceNameListKey(common.BuildJobPrefix)
)

// 添加一个请求到当前实例的缓存列表中
func AddRequest(m *models.Request) error {
	requestKey := GenRequestKey(common.BuildJobPrefix, m.InstanceName)
	requestJsonData, err := json.Marshal(m)
	if err != nil {
		log.Error("ERROR: ", err)
		return err
	}
	err = RedisClient.HSet(requestKey, GenFieldByRequest(m), requestJsonData).Err()
	if err != nil {
		return err
	}
	return nil
}

func UpdateRequest(m *models.Request) error {
	return AddRequest(m)
}

func DeleteRequest(m *models.Request) error{
	requestKey := GenRequestKey(common.BuildJobPrefix, m.InstanceName)
	return RedisClient.HDel(requestKey, GenFieldByRequest(m)).Err()
}

// 查询某个特定的请求，需要从全局查询
func GetRequestByNameAndRequestType(name string, requestType string) (*models.Request, error) {
	instanceNameList, err := GetInstanceNameList()
	if err != nil {
		return nil, err
	}
	var request *models.Request
	for _, instanceName := range instanceNameList {
		request, err = GetRequestByNameAndRequestTypeAndInstanceName(name, requestType, instanceName)
		if err == nil {
			return request, nil
		}
	}
	return nil, err
}

func GetAllRequestByInstanceName(instanceName string) ([]*models.Request, error) {
	requestKey2 := GenRequestKey(common.BuildJobPrefix, instanceName)
	requestMap, err := RedisClient.HGetAll(requestKey2).Result()
	if err != nil {
		return nil, err
	}
	requestList := make([]*models.Request, len(requestMap))
	i := 0
	for _, requestJsonData := range requestMap {
		m := &models.Request{}
		err = json.Unmarshal([]byte(requestJsonData), m)
		if err != nil {
			return nil,err
		}
		requestList[i] = m
		i += 1
	}
	return requestList, nil
}

func GetRequestByNameAndRequestTypeAndInstanceName(name string, requestType string, instanceName string) (*models.Request, error){
	requestKey2 := GenRequestKey(common.BuildJobPrefix, instanceName)
	requestJsonData, err := RedisClient.HGet(requestKey2, GenFieldByRequestTypeAndName(requestType, name)).Result()
	if err != nil {
		return nil, err
	}
	m := &models.Request{}
	err = json.Unmarshal([]byte(requestJsonData), m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func GetInstanceNameList() ([]string, error) {
	instanceNameListJsonData, err := RedisClient.Get(instanceNameListKey).Result()
	instanceNameList := make([]string, 0)
	// 不存在，不用操作
	if err == redis.Nil {
		return instanceNameList, nil
	} else if err != nil {
		return nil, err
	} else {
		// 存在
		err = json.Unmarshal([]byte(instanceNameListJsonData), &instanceNameList)
		if err != nil {
			return nil, err
		}
	}
	return instanceNameList, nil
}



// redis key键设置
func GenMetaDistributeKey(prefix string) string {
	return fmt.Sprintf("%s/%s", prefix, "meta")
}

func GenInstanceNameListKey(prefix string) string {
	return fmt.Sprintf("%s/%s", prefix, "instance-name-list")
}

func GenRequestKey(prefix string, instanceName string) string {
	return fmt.Sprintf("%s/%s/%s", prefix, instanceName, "requests")
}

func GenInstanceKey(prefix string, instanceName string) string {
	return fmt.Sprintf("%s/%s/%s", prefix, "instances", instanceName)
}

func GenFieldByRequest(r *models.Request) string {
	return GenFieldByRequestTypeAndName(r.RequestType, r.Name)
}

func GenFieldByRequestTypeAndName(requestType string, name string) string {
	return fmt.Sprintf("%s/%s", requestType, name)
}