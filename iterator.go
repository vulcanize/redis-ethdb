package redisethdb

import (
	"context"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/go-redis/redis"
)

var _ ethdb.Iterator = &Iterator{}

var (
	defaultScanSize int64 = 1024
)

// Iterating redis interface required for supporting ethdb.Iterator
type Iterating interface {
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
	io.Closer
}

// Iterator implements ethdb.Iterator on top of Redis using scans
// it is not safe for concurrent use, but multiple iterators can act concurrently on the same DB
type Iterator struct {
	// iterating interface
	it Iterating

	// iterator state
	prefix          string
	init            bool
	ctx             context.Context
	scanSize        int64         // the number of results to return per set
	cursor          uint64        // the current cursor location
	currentKeys     []string      // the current set of keys
	currentValues   []interface{} // the current set of values
	currentSetIndex uint          // the index in the currentSet
	err             error
}

// NewIterator creates a new Iterator
func NewIterator(ctx context.Context, prefix string, scanSize int64, it Iterating) ethdb.Iterator {
	if scanSize == 0 {
		scanSize = defaultScanSize
	}
	return &Iterator{
		prefix:   prefix,
		init:     false,
		ctx:      ctx,
		it:       it,
		scanSize: scanSize,
	}
}

// Next satisfies ethdb.Iterator
func (i *Iterator) Next() bool {
	// if the cursor is at 0 after init then we have finished iteration
	if i.cursor == 0 && i.init {
		return false
	}

	// local state is sufficient
	if int(i.currentSetIndex) < len(i.currentKeys)-1 {
		i.currentSetIndex++
		return true
	}

	// we need to retrieve and begin iterating the next set
	keys, cursor, err := i.it.Scan(i.ctx, i.cursor, i.prefix, i.scanSize).Result()
	if err != nil {
		i.err = fmt.Errorf("it.Next(): Scan() error: %v", err)
		return false
	}
	vals, err := i.it.MGet(i.ctx, keys...).Result()
	if err != nil {
		i.err = fmt.Errorf("it.Next(): MGet() error: %v", err)
		return false
	}

	// this shouldn't be necessary as the redis client should guarantee this, but let's put this guard here for now
	if len(vals) != len(keys) {
		i.err = fmt.Errorf("number of values must match the number of keys")
		return false
	}

	i.currentKeys = keys
	i.currentValues = vals
	i.currentSetIndex = 0
	i.cursor = cursor
	i.init = true

	return true
}

// Error satisfies ethdb.Iterator
func (i *Iterator) Error() error {
	return i.err
}

// Key satisfies ethdb.Iterator
func (i *Iterator) Key() []byte {
	return []byte(i.currentKeys[i.currentSetIndex])
}

// Value satisfies ethdb.Iterator
func (i *Iterator) Value() []byte {
	return []byte(i.currentValues[i.currentSetIndex].(string))
}

// Release satisfies ethdb.Iterator
func (i *Iterator) Release() {
	if err := i.it.Close(); err != nil {
		i.err = err
		return
	}
	i.currentSetIndex = 0
	i.cursor = 0
	i.currentValues = nil
	i.currentKeys = nil
	i.err = nil
}
