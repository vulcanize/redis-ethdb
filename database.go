package redisethdb

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/ethereum/go-ethereum/ethdb"
)

var _ ethdb.Database = &Database{}

// Redis interface required for supporting ethdb interfaces
type Redis interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
}

type Database struct {
	redis Redis
}

func NewDatabase(r Redis) *Database {
	return &Database{
		redis: r,
	}
}

func NewClientDatabase(cl *redis.Client) *Database {
	return &Database{
		redis: cl,
	}
}

func NewClusterDatabase(cl *redis.ClusterClient) *Database {
	return &Database{
		redis: cl,
	}
}

// Has satisfies ethdb.KeyValueReader
func (db *Database) Has(key []byte) (bool, error) {
	panic("implement me")
}

// Get satisfies ethdb.KeyValueReader
func (db *Database) Get(keu []byte) ([]byte, error) {
	panic("implement me")
}

// Put satisfies ethdb.KeyValueWriter
func (db *Database) Put(key, value []byte) error {
	panic("implement me")
}

// Delete satisfies ethdb.KeyValueWriter
func (db *Database) Delete(key []byte) error {
	panic("implement me")
}

// Stat satisfied ethdb.Stater
func (db *Database) Stat(property string) (string, error) {
	panic("implement me")
}

// Compact satisfied ethdb.Compacter
func (db *Database) Compact(start []byte, limit []byte) error {
	panic("implement me")
}

// HasAncient satisfied ethdb.AncientReader
func (db *Database) HasAncient(kind string, number uint64) (bool, error) {
	panic("implement me")
}

// Ancient satisfied ethdb.AncientReader
func (db *Database) Ancient(kind string, number uint64) ([]byte, error) {
	panic("implement me")
}

// ReadAncients satisfied ethdb.AncientReader
func (db *Database) ReadAncients(kind string, start, count, maxBytes uint64) ([][]byte, error) {
	panic("implement me")
}

// Ancients satisfied ethdb.AncientReader
func (db *Database) Ancients() (uint64, error) {
	panic("implement me")
}

// AncientSize satisfied ethdb.AncientReader
func (db *Database) AncientSize(kind string) (uint64, error) {
	panic("implement me")
}

// ModifyAncients satisfied ethdb.AncientWriter
func (db *Database) ModifyAncients(func(ethdb.AncientWriteOp) error) (int64, error) {
	panic("implement me")
}

// TruncateAncients satisfied ethdb.AncientWriter
func (db *Database) TruncateAncients(n uint64) error {
	panic("implement me")
}

// Sync satisfied ethdb.AncientWriter
func (db *Database) Sync() error {
	panic("implement me")
}

// Close satisfied io.Closer
func (db *Database) Close() error {
	panic("implement me")
}
