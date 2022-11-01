/**
 * etcd实现分布式锁
 *
 * @title etcd_lock
 * @projectName distributeLock
 * @author sauryniu
 * @date 2022/10/24 16:41
 */

package distributelock

import (
	"context"
	"errors"
	"fmt"
	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync/atomic"
	"time"
)

const (
	lockType    = 1
	tryLockType = 2
)

var notInitErr = errors.New("etcd locker: The locker was not initialized")

type etcdLock struct {
	addr     string
	initFlag atomic.Bool
	client   *v3.Client
	ttl      int
}

func (e *etcdLock) Lock(ctx context.Context, key string, ops ...OpOption) (Unlocker, error) {
	op := &Op{ttl: defaultWaitSeconds}
	for _, opt := range ops {
		opt(op)
	}

	var cancel context.CancelFunc
	if op.ttl > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(op.ttl))
	}
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	return e.doLock(ctx, key, lockType)
}

func (e *etcdLock) TryLock(ctx context.Context, key string) (Unlocker, error) {
	return e.doLock(ctx, key, tryLockType)
}

func (e *etcdLock) doLock(ctx context.Context, key string, t int) (Unlocker, error) {
	if !e.isInit() {
		return nil, notInitErr
	}

	session, err := concurrency.NewSession(e.client, concurrency.WithTTL(e.ttl))
	if err != nil {
		return nil, err
	}

	lockKey := fmt.Sprintf("/dLock/%s", key)
	mutex := concurrency.NewMutex(session, lockKey)

	if t == lockType {
		err = mutex.Lock(ctx)
	} else {
		err = mutex.TryLock(ctx)
	}

	if err != nil {
		return nil, err
	}
	unlocker := func(ctx2 context.Context) error {
		return mutex.Unlock(ctx)
	}
	return unlocker, nil
}

func (e *etcdLock) init() {
	var err error
	e.client, err = v3.New(v3.Config{
		Endpoints:   []string{e.addr},
		DialTimeout: time.Second * 10,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	e.initFlag.Store(true)
}

func (e *etcdLock) isInit() bool {
	if !e.initFlag.Load() {
		e.init()
	}
	return e.initFlag.Load()
}

func newEtcdLock(serverAddr string, ttl int) *etcdLock {
	lock := &etcdLock{
		addr: serverAddr,
		ttl:  ttl,
	}
	lock.init()
	return lock
}
