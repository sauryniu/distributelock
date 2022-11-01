/**
 * 分布式锁接口
 *
 * @title distribute_lock
 * @projectName distributeLock
 * @author sauryniu
 * @date 2022/10/24 14:19
 */

package distributelock

import (
	"context"
)

type Unlocker func(ctx context.Context) error

type LockType int

const (
	EtcdLock           LockType = 0
	defaultWaitSeconds          = -1
)

type Op struct {
	ttl int
}

type OpOption func(*Op)

func WithTTL(ttl int) OpOption {
	return func(op *Op) {
		if ttl > 0 {
			op.ttl = ttl
		}
	}
}

type Locker interface {
	TryLock(ctx context.Context, key string) (Unlocker, error)
	Lock(ctx context.Context, key string, ops ...OpOption) (Unlocker, error)
}

func NewLocker(serverAddr string, ttl int, lockType LockType) Locker {
	switch lockType {
	case EtcdLock:
		return newEtcdLock(serverAddr, ttl)
	}
	return nil
}
