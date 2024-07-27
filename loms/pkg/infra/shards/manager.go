package shards

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spaolacci/murmur3"
)

var ErrShardIndexOutOfRange = errors.New("shard index out of range")

type ShardIndex int
type ShardKey string
type ShardFn func(ShardKey) ShardIndex

type Manager struct {
	fn   ShardFn
	pool []*pgxpool.Pool
}

func NewManager(fn ShardFn, pools []*pgxpool.Pool) *Manager {
	return &Manager{
		fn:   fn,
		pool: pools,
	}
}

func (m *Manager) GetShardIndex(key ShardKey) ShardIndex {
	return m.fn(key)
}

func (m *Manager) GetShardIndexFromID(key int) ShardIndex {
	return ShardIndex(key % 10)
}

func (m *Manager) GetShardByIndex(index ShardIndex) (*pgxpool.Pool, error) {
	if int(index) >= len(m.pool) {
		return nil, ErrShardIndexOutOfRange
	}

	return m.pool[index], nil
}

func (m *Manager) GetShards() []*pgxpool.Pool {
	return m.pool
}

func GetMurmur3ShardFn(shardsCount int) ShardFn {
	hasher := murmur3.New32()
	return func(key ShardKey) ShardIndex {
		defer hasher.Reset()

		_, _ = hasher.Write([]byte(key))

		return ShardIndex(hasher.Sum32() % uint32(shardsCount))
	}
}
