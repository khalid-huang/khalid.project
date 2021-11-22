package async

import (
	"bryson.foundation/kbuildresource/cache"
	"bryson.foundation/kbuildresource/models"
	"github.com/sirupsen/logrus"
)

type RequestHandler interface {
	// 在请求进入执行队列之前，做一些前置处理，诸如参数校验，默认值填充等工作，values里面保存了在处理链上需要传递的一些中间信息，避免重复计算
	// 这一步，无论是同步还是异步都需要做，是一些比较轻量级的操作
	PreExec(requestDTO interface{}, requestType string, values map[string]interface{}) error
	// 根据请求参数封装出request，并将request进行缓存； 记得把dto里面的instanceName注入到request中
	MakeRequest(requestDTO interface{}, requestType string, values map[string]interface{}) (*models.Request, error)
	// 异步执行请求，limitChan用于控制并发量
	AsyncExec(request *models.Request, limitChan <-chan struct{})
	// 在请求进入执行队列之后需要做的一些操作
	PostAsyncExec(request *models.Request, requestType string, values map[string]interface{}) error
	// 同步执行请求
	SyncExec(requestDTO interface{}, requestType string, values map[string]interface{}) (interface{}, error)
	// 给requestDTO 设置请求的instanceName
	SetInstanceName(requestDTO interface{}, instanceName string)
	// 用于生成异步执行的Response,返回给客户端，用于处理信息，隐藏一些没有必要返回给客户端的信息，同步的场景会直接处理返回值
	MakeAsyncResponse(requestDTO interface{}, requestType string, values map[string]interface{}) interface{}
	// 用于处理从deadInstance中接管request请求
	HandleTakeOverRequest(request *models.Request, newInstanceName string) error
}

// 接收Pending request时需要做的处理，这是通用广场
func HandleCacheDataForTakeOverPendingRequest(request *models.Request, newInstanceName string) error {
	// 在原本的redis hash 里面进行删除， 在新的列表里面添加
	err := cache.DeleteRequest(request)
	if err != nil {
		logrus.Error("ERROR: take over job failed, error: ", err)
		return err
	}
	// 设置新的instanceName
	request.InstanceName = newInstanceName
	err = cache.AddRequest(request)
	if err != nil {
		logrus.Error("ERROR: take over job failed, error: ", err)
		return err
	}
	return nil
}