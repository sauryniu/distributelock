/**
 * 测试文件
 *
 * @title distribute_lock_test.go
 * @projectName distributeLock
 * @author sauryniu
 * @date 2022/10/31 11:17
 */

package distributelock

import (
	"context"
	"fmt"
	"testing"
)

func TestNewDistributeLock(t *testing.T) {
	type args struct {
		serverAddr string
		ttl        int
	}
	type args2 struct {
		key string
	}
	tests := []struct {
		name  string
		args  args
		args2 []args2
		want  DistributeLock
	}{
		{
			name: "1",
			args: args{
				serverAddr: "10.1.30.79:12379",
				ttl:        30,
			},
			args2: []args2{
				{key: "111"},
				{key: "222"},
				{key: "333"},
				{key: "444"},
				{key: "555"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locker := NewDistributeLock(tt.args.serverAddr, tt.args.ttl)
			for _, arg := range tt.args2 {
				fmt.Println(arg.key)
				unlocker, err := locker.Lock(context.Background(), arg.key)
				if err != nil {
					panic(err)
				}
				if err = unlocker(context.Background()); err != nil {
					panic(err)
				}
			}

		})
	}
}
