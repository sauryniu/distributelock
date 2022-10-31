/**
 * 分布式锁接口
 *
 * @title distribute_lock
 * @projectName distributeLock
 * @author sauryniu
 * @date 2022/10/24 14:19
 */

package distributelock

import "context"

type Unlocker func(ctx context.Context) error

type LockType int

const (
	EtcdLock LockType = 0
)

type DistributeLock interface {
	Lock(ctx context.Context, key string) (Unlocker, error)
}

func NewDistributeLock(serverAddr string, ttl int, lockType LockType) DistributeLock {
	switch lockType {
	case EtcdLock:
		return newEtcdLock(serverAddr, ttl)
	}
	return nil
}
