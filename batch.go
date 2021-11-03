package redisethdb

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/go-redis/redis"
)

var _ ethdb.Batch = &Batch{}

// Batching redis interface required for supporting ethdb.Batch
type Batching interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Discard() error
	Exec(ctx context.Context) ([]redis.Cmder, error)
}

// Batch implements ethdb.Batch on top of Redis using pipelined transactions
type Batch struct {
	// batching interface
	pipeline Batching

	// batch state
	ctx         context.Context
	valueSize   *big.Int // size in bytes
	replayCache map[string][]byte
}

// NewBatch creates a new Batch
func NewBatch(ctx context.Context, b Batching) ethdb.Batch {
	return &Batch{
		ctx:         ctx,
		pipeline:    b,
		valueSize:   big.NewInt(0),
		replayCache: make(map[string][]byte),
	}
}

// Put satisfies ethdb.Batch
func (b Batch) Put(key []byte, value []byte) error {
	byteSize := int64(len(value))
	strKey := common.Bytes2Hex(key)
	b.valueSize.Add(b.valueSize, big.NewInt(byteSize))
	b.replayCache[strKey] = value
	b.pipeline.Set(b.ctx, strKey, value, 0)
	return nil
}

// Delete satisfies ethdb.Batch
func (b Batch) Delete(key []byte) error {
	b.pipeline.Del(b.ctx, common.Bytes2Hex(key))
	return nil
}

// ValueSize satisfies ethdb.Batch
// it returns it in units of GB
func (b Batch) ValueSize() int {
	gb := new(big.Int).Div(b.valueSize, big.NewInt(1000000))
	return int(gb.Int64())
}

// Write satisfies ethdb.Batch
func (b Batch) Write() error {
	res, err := b.pipeline.Exec(b.ctx)
	if err == nil {
		b.replayCache = nil
		return nil
	}
	var errMsg string
	if err != nil {
		for _, r := range res {
			errMsg += r.Name() + " : " + r.Err().Error() + "\r\n"
		}
	}
	return fmt.Errorf(err.Error() + "\r\n" + errMsg)
}

// Reset satisfies ethdb.Batch
func (b Batch) Reset() {
	b.pipeline.Discard()
	b.replayCache = make(map[string][]byte)
	b.valueSize.SetUint64(0)
}

// Replay satisfies ethdb.Batch
func (b Batch) Replay(w ethdb.KeyValueWriter) error {
	for key, value := range b.replayCache {
		if err := w.Put(common.Hex2Bytes(key), value); err != nil {
			return err
		}
	}
	b.replayCache = make(map[string][]byte)
	return nil
}
