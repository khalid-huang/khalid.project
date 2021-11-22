package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"time"
)

// 通过redis获取分布式锁
// @Param lockName 锁的名字
// @Param acquireTimeout  等待的时长
// @Param lockTimeout 锁的有效期
func RedisGetLock(lockName string, acquireTimeout time.Duration, lockTimeout time.Duration) (bool, error) {
	code := time.Now().String()
	endTime := time.Now().Add(acquireTimeout).UnixNano()
	for time.Now().UnixNano() <= endTime {
		success, err := RedisClient.SetNX(lockName, code, lockTimeout).Result()
		if err != nil {
			return false, err
		} else if success {
			return true, nil
		} else if RedisClient.TTL(lockName).Val() == -1 {
			RedisClient.Expire(lockName, lockTimeout)
		} else if !success {
			//锁已被占用
			return false, nil
		}
		time.Sleep(time.Millisecond)
	}
	return false, fmt.Errorf("timeout")
}

func RedisReleaseLock(lockName string) bool {
	_, err := RedisClient.Del(lockName).Result()
	if err != nil && err != redis.Nil {
		return false
	}
	return true
}

func IsExpire(key string) bool {
	// 过期或不存在都返回这个
	// TTL 方法，不存在返回-2， 过期返回-1
	return RedisClient.TTL(key).Val() == -2 * time.Second || RedisClient.TTL(key).Val() == -1 * time.Second
}

// 获取锁，非阻塞
func LockKey(key string, lockLeaseTime time.Duration) (bool, error) {
	success, err := RedisClient.SetNX(key, "1", lockLeaseTime).Result()
	if err != nil {
		return success, err
	}
	logrus.Infof("INFO: lock key %s success", key)
	return success, nil
}

func RenewExpiration(key string, lockLeaseTime time.Duration) error {
	success, err := RedisClient.Expire(key, lockLeaseTime).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	// 不存在或者没有续期成功
	if !success && IsExpire(key) {
		log.Info("INFO: instance %s retry get lock", key)
		_, err := LockKey(key, lockLeaseTime)
		if err != nil {
			return err
		}
		success = true
	}
	return nil
}

func DelKey(key string) error {
	_, err := RedisClient.Del(key).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}