package redisethdb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/go-redis/redis/v8"
)

var _ ethdb.Database = &Database{}
var _ Redis = &redis.Client{}
var _ Redis = &redis.ClusterClient{}
var _ Redis = &redis.Ring{}

var (
	errNotSupported = errors.New("this operation is not supported")
)

// Redis interface required for supporting ethdb interfaces
type Redis interface {
	// Database interfaces
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd

	// Iterator interfaces
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd

	// Batch interfaces
	TxPipeline() redis.Pipeliner

	// Stats interfaces
	PoolStats() *redis.PoolStats
	DBSize(ctx context.Context) *redis.IntCmd
	Info(ctx context.Context, section ...string) *redis.StringCmd

	io.Closer
}

// Database implements ethdb.Database on top of Redis
type Database struct {
	ctx   context.Context
	redis Redis
}

// NewDatabase creates a new Database
func NewDatabase(ctx context.Context, r Redis) ethdb.Database {
	return &Database{
		ctx:   ctx,
		redis: r,
	}
}

// NewClientDatabase creates a redis-ethdb using a Client
func NewClientDatabase(ctx context.Context, cl *redis.Client) ethdb.Database {
	return &Database{
		ctx:   ctx,
		redis: cl,
	}
}

// NewClusterDatabase creates a redis-ethdb using a ClusterClient
func NewClusterDatabase(ctx context.Context, cl *redis.ClusterClient) ethdb.Database {
	return &Database{
		ctx:   ctx,
		redis: cl,
	}
}

// NewRingDatabase creates a redis-ethdb using a Ring
func NewRingDatabase(ctx context.Context, cl *redis.Ring) ethdb.Database {
	return &Database{
		ctx:   ctx,
		redis: cl,
	}
}

// Has satisfies ethdb.KeyValueReader
func (db *Database) Has(key []byte) (bool, error) {
	res, err := db.redis.Exists(db.ctx, common.Bytes2Hex(key)).Result()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// Get satisfies ethdb.KeyValueReader
func (db *Database) Get(key []byte) ([]byte, error) {
	res, err := db.redis.Get(db.ctx, common.Bytes2Hex(key)).Result()
	if err != nil {
		return nil, err
	}
	return []byte(res), err
}

// Put satisfies ethdb.KeyValueWriter
func (db *Database) Put(key, value []byte) error {
	return db.redis.Set(db.ctx, common.Bytes2Hex(key), value, 0).Err()
}

// Delete satisfies ethdb.KeyValueWriter
func (db *Database) Delete(key []byte) error {
	return db.redis.Del(db.ctx, common.Bytes2Hex(key)).Err()
}

// Stat satisfies ethdb.Stater
func (db *Database) Stat(property string) (string, error) {
	if inList(Stat(property), poolStats) {
		poolInfo := db.redis.PoolStats()
		switch property {
		case HITS:
			return strconv.Itoa(int(poolInfo.Hits)), nil
		case MISSES:
			return strconv.Itoa(int(poolInfo.Misses)), nil
		case TIMEOUTS:
			return strconv.Itoa(int(poolInfo.Timeouts)), nil
		case TOTAL_CONNS:
			return strconv.Itoa(int(poolInfo.TotalConns)), nil
		case IDLE_CONNS:
			return strconv.Itoa(int(poolInfo.IdleConns)), nil
		case STABLE_CONNS:
			return strconv.Itoa(int(poolInfo.StaleConns)), nil
		default:
			return "", fmt.Errorf("unrecognized PoolStats property: %s", property)
		}
	}

	if inList(Stat(property), infoStats) {
		return db.redis.Info(db.ctx, property).Result()
	}

	if inList(Stat(property), dataStats) {
		switch property {
		case DB_SIZE:
			size, err := db.redis.DBSize(db.ctx).Result()
			if err != nil {
				return "", err
			}
			return strconv.Itoa(int(size)), nil
		default:
			return "", fmt.Errorf("unrecognized DataStats property: %s", property)
		}
	}

	return "", fmt.Errorf("unrecognized property: %s", property)
}

// Compact satisfies ethdb.Compacter
func (db *Database) Compact([]byte, []byte) error {
	return errNotSupported
}

// HasAncient satisfies ethdb.AncientReader
func (db *Database) HasAncient(string, uint64) (bool, error) {
	return false, errNotSupported
}

// Ancient satisfies ethdb.AncientReader
func (db *Database) Ancient(string, uint64) ([]byte, error) {
	return nil, errNotSupported
}

// ReadAncients satisfies ethdb.AncientReader
func (db *Database) ReadAncients(string, uint64, uint64, uint64) ([][]byte, error) {
	return nil, errNotSupported
}

// Ancients satisfies ethdb.AncientReader
func (db *Database) Ancients() (uint64, error) {
	return 0, errNotSupported
}

// AncientSize satisfies ethdb.AncientReader
func (db *Database) AncientSize(string) (uint64, error) {
	return 0, errNotSupported
}

// ModifyAncients satisfies ethdb.AncientWriter
func (db *Database) ModifyAncients(func(ethdb.AncientWriteOp) error) (int64, error) {
	return 0, errNotSupported
}

// TruncateAncients satisfies ethdb.AncientWriter
func (db *Database) TruncateAncients(uint64) error {
	return errNotSupported
}

// Sync satisfies ethdb.AncientWriter
func (db *Database) Sync() error {
	return errNotSupported
}

// Close satisfies io.Closer
func (db *Database) Close() error {
	return db.redis.Close()
}

// NewIterator satisfies ethdb.Iteratee
// Note: there is no practical way to have Redis begin iteration at an arbitrary "start" key, so this arg is discarded
func (db *Database) NewIterator(prefix, _ []byte) ethdb.Iterator {
	return NewIterator(db.ctx, common.Bytes2Hex(prefix), defaultScanSize, db.redis)
}

// NewBatch satisfies ethdb.Batcher
func (db *Database) NewBatch() ethdb.Batch {
	return NewBatch(db.ctx, db.redis.TxPipeline())
}
