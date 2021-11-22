package instance

import (
	"bryson.foundation/kbuildresource/cache"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	intervalLockLeaseTime  = 6 *time.Second // 每次续期时长
	RenewInterval = 3 * time.Second // 每3秒续期一次
	commonRetryTimes = 3
	commonRetryInterval = 100
	retryAccessMasterInterval = 15
	ROLEMASTER = "master"
	ROLEBACKUP = "backup"
)

// 场景是跟随实例应用一同，实例应用活着的时候，如果一直是master就一直master，不断续租，比如多个实例时，只需要一个实例进行list-watch就可以了
// 用于实现一种主备的机制实现，传入两种回调函数，一种是成为master时的函数，一种是失败master身份的函数，以及一个任务标识
// 默认状态都是backup
type masterCallbackFunc func()  // 从backup变成master时要调用的函数
type backupCallbackFunc func()  // 从master变成backup时要调用的函数

type MasterBackupJob struct {
	masterCallback masterCallbackFunc
	backupCallback backupCallbackFunc
	name           string // 任务名称，唯一标识，也做为redis分布式锁的标识
	signalCh chan os.Signal // 监听系统信号
	role string // 当前角色，默认backup
}

// 都要实现优雅停机的机制
func (m *MasterBackupJob) StartUp() {
	signal.Notify(m.signalCh, syscall.SIGINT, syscall.SIGTERM)
	t := time.NewTicker(retryAccessMasterInterval * time.Second)
	for {
		select {
		case <-t.C:
			if m.role != ROLEMASTER {
				m.tryBecomeMaster()
			} else {
				logrus.Infof("I am master of MasterBackupJob %s", m.name)
			}
		case signalVal := <-m.signalCh:
			logrus.Infof("INFO: shutting down the job because signal %s", signalVal)
			m.preStop()
		}
	}
}

// 尝试获取锁成为master，执行任务
func (m *MasterBackupJob) tryBecomeMaster() {
	var result bool
	var err error
	for i := 0; i < commonRetryTimes; i++ {
		result, err = cache.LockKey(m.name, intervalLockLeaseTime)
		if err == nil {
			break
		} else {
			logrus.Error("ERROR: get lock failed, err is ", err)
		}
		time.Sleep(commonRetryInterval * time.Millisecond)
	}
	if result {
		m.role = ROLEMASTER
		logrus.Infof("INFO: become master, start job %s", m.name)
		go m.masterCallback()
		// 不断进行续期，直到无法续期
		ticker := time.NewTicker(RenewInterval)
		go func() {
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					var err error
					for i := 0; i < commonRetryTimes; i++ {
						err := cache.RenewExpiration(m.name, intervalLockLeaseTime)
						if err == nil {
							logrus.Infof("INFO: RenewExpiration success")
							break
						} else {
							logrus.Error("ERROR: renew lock failed")
						}
					}
					if err != nil {
						m.role = ROLEBACKUP
						logrus.Info("INFO: renewexpiration faield, become worker")
						go m.backupCallback()
					}
				}
			}
		}()
	} else {
		m.role = ROLEBACKUP
		logrus.Info("INFO: become worker, no need to start job %s", m.name)
	}
}

func (m *MasterBackupJob) preStop() {
	// 如果是master就调用
	if m.role == ROLEMASTER {
		go m.backupCallback()
	}
}
