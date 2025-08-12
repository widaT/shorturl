package shint

import (
	"context"
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

type Picker func(size uint64) uint64

type Shint struct {
	Shard       uint64
	ShardSubfix int
	Picker      Picker
	Key         string
	client      *redis.Client
}

func New(shard uint64, key string,
	client *redis.Client) *Shint {
	sh := new(Shint)
	sh.Shard = CeilPowerOf2(shard)
	sh.ShardSubfix = int(math.Log2(float64(sh.Shard)))
	sh.Picker = RRPicker()
	sh.Key = key
	sh.client = client
	return sh
}

func CeilPowerOf2(n uint64) uint64 {
	if n <= 1 {
		return 1
	}
	// bits.Len(x) 返回表示x所需的位数。
	// 通过对 n-1 使用 Len，可以巧妙地处理n本身就是2的幂的情况。
	return 1 << bits.Len64(n-1)
}

func RRPicker() Picker {
	var n atomic.Uint64
	return func(size uint64) uint64 {
		mod := size - 1
		return n.Add(1) & mod
	}
}

func (sh *Shint) Incr(ctx context.Context) (uint64, error) {
	shard := sh.Picker(sh.Shard)
	fmt.Println(shard)
	v, err := sh.incrby(ctx, shard)
	if err != nil {
		return 0, err
	}
	return sh.shift(v, shard), nil
}

func (sh *Shint) shift(v, shard uint64) uint64 {
	return v<<sh.ShardSubfix | shard
}

func (sh *Shint) incrby(ctx context.Context, shard uint64) (uint64, error) {
	key := sh.Key + ":" + strconv.FormatUint(shard, 10)
	fmt.Println(key)
	v, err := sh.client.IncrBy(ctx, key, 1).Result()
	if err != nil {
		return 0, err
	}
	return uint64(v), nil
}
