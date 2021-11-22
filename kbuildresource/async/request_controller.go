package async

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"

	"bryson.foundation/kbuildresource/cache"
	"bryson.foundation/kbuildresource/common"
	"bryson.foundation/kbuildresource/models"
)

type RequestController struct {
	limitChan chan struct{} //用于控制并发，可以用协程池来做
	stopCh chan struct{} // 控制器停止通道
	requestChannel chan *models.Request
	instanceName string // 对应的实例的名字
	stopping int // 表明是否是正在停止中
}

var r *RequestController
var requestHandlerMap = make(map[string]RequestHandler, 0) //保存了各种请求类型的处理器，实现依赖反转，避免每次新生成一种处理器，都要修改函数

func NewRequestController(instanceName string) *RequestController {
	r = &RequestController{
		limitChan:      make(chan struct{}, 100), // 这个要控制小点，避免造成数据库连接过多
		stopCh:         make(chan struct{}),
		requestChannel: make(chan *models.Request, 2000),
		stopping: 0,
		instanceName: instanceName,
	}
	return r
}

func GetRequestController() *RequestController {
	return r
}

// 启动请求控制器
func (r *RequestController) StartUp() {
	for request := range r.requestChannel {
		if r.stopping == 1 {
			log.Info("INFO: request controller is stopping skip exec request")
			continue
		}
		r.limitChan <- struct{}{}
		logrus.Infof("INFO: receive request %s and start handle", request.Name)
		requestHandler,_ := getHandlerFromRequestType(request.RequestType)
		go requestHandler.AsyncExec(request, r.limitChan)
	}
	logrus.Info("INFO: finish requestChannel")
	close(r.stopCh) // 通知shutdown函数继续执行
}

func (r *RequestController) Shutdown() {
	logrus.Info("INFO: shutdown requestController")
	// sleep 一小段时间，保证收到的请求都入channel了
	time.Sleep(2 * time.Second)
	r.stopping = 1 // 表明正在关闭
	close(r.requestChannel) // 关闭requestChannel,促使requestHandler里面的for range循环可以在遍历完成之后结束
	<-r.stopCh // 等待requestHandle处理完成的信号，当close(r.stopCh)时可以结束
}

// 请求管理器对外提供的接收请求的接口
// @Param requestDTO interface{} 请求传输对象，主要包含用户传入的参数
// @Param requestType string 请求类型，用于分派请求到对应的处理器
// return interface{} 请求处理的返回结果
func (r *RequestController) AcceptRequest(requestDTO interface{}, requestType string) (interface{}, error) {
	requestHandler, _ := getHandlerFromRequestType(requestType)
	values := make(map[string]interface{}, 0)
	requestHandler.SetInstanceName(requestDTO, r.instanceName)
	err := requestHandler.PreExec(requestDTO, requestType, values)
	if err != nil {
		return requestDTO, err
	}
	request, err := requestHandler.MakeRequest(requestDTO, requestType, values)
	if err != nil {
		logrus.Error("ERROR: MakeRequest failed, try to use SyncExec")
		return requestHandler.SyncExec(requestDTO, requestType, values)
	}
	go r.sendRequestToChannel(request)
	err = requestHandler.PostAsyncExec(request, requestType, values)
	if err != nil {
		logrus.Error("ERROR: posyAsyncExec failed, err: ", err)
	}
	return requestDTO, nil
}

func (r *RequestController) sendRequestToChannel(request *models.Request) {
	r.requestChannel <- request
}

func (r *RequestController) TakeOverRequest(deadInstanceName string, newInstanceName string) error {
	logrus.Infof("INFO: start takeover request of instance %s", deadInstanceName)
	requests, err := cache.GetAllRequestByInstanceName(deadInstanceName)
	if err != nil {
		return err
	}
	for _, request := range requests {
		logrus.Infof("INFO: take over instance-%s request %s which status is %s", deadInstanceName, request.Name, request.Status)
		// 如果重启状态为pending，直接放入队列中，等待执行
		requestHandler, err := getHandlerFromRequestType(request.RequestType)
		if err != nil {
			return err
		}
		err = requestHandler.HandleTakeOverRequest(request, r.instanceName)
		if err != nil {
			return err
		}
		go r.sendRequestToChannel(request)
	}
	return nil
}

func handleCacheDataForTakeOver(request *models.Request, newInstanceName string) error {
	// 在原本的instance里面删除，再加入到新的里面
	err := cache.DeleteRequest(request)
	if err != nil {
		return err
	}
	request.InstanceName = newInstanceName
	request.Status = common.RequestStatusPending // 这里简单假设，所有的都是需要重新执行的
	err = cache.AddRequest(request)
	if err != nil {
		return err
	}
	return nil
}

func RegisterRequestHandler(requestType string, requestHandler RequestHandler) {
	requestHandlerMap[requestType] = requestHandler
}

func getHandlerFromRequestType(requestType string) (RequestHandler, error) {
	s := strings.Split(requestType, "_")
	// 使用依赖反转进行替换
	//switch s[0] {
	//case common.BuildJobPrefix:
	//	return handler.GetBuildJobHandler()
	//default:
	//	return nil
	//}
	requestHandler, ok := requestHandlerMap[s[0]]
	if !ok {
		return nil, fmt.Errorf("invalid requestType %s", s[0])
	}
	return requestHandler, nil
}