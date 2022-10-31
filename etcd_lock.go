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
	"crypto/rand"
	"errors"
	"fmt"
	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"math"
	"math/big"
	"sync/atomic"
	"time"
)

type etcdLock struct {
	addr    string
	isInit  atomic.Bool
	client  *v3.Client
	session *concurrency.Session
	ttl     int
}

func (e *etcdLock) Lock(ctx context.Context, key string) (Unlocker, error) {
	if !e.isInit.Load() {
		e.init()
		return nil, errors.New("not init")
	}
	r, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		return nil, err
	}
	prefix := fmt.Sprintf("/dLock/%s/%d", key, r)
	mutex := concurrency.NewMutex(e.session, prefix)
	err = mutex.Lock(ctx)
	if err != nil {
		return nil, err
	}
	unlocker := func(ctx2 context.Context) error {
		return mutex.Unlock(ctx2)
	}
	return unlocker, nil
}

func (e *etcdLock) init() {
	var err error
	e.client, err = v3.New(v3.Config{
		Endpoints:   []string{"10.1.30.79:12379"},
		DialTimeout: time.Second * 10,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	e.session, err = concurrency.NewSession(e.client, concurrency.WithTTL(e.ttl))
	if err != nil {
		fmt.Println(err)
		return
	}

	e.isInit.Store(true)
}

func newEtcdLock(serverAddr string, ttl int) *etcdLock {
	lock := &etcdLock{
		addr: serverAddr,
		ttl:  ttl,
	}
	lock.init()
	return lock
}