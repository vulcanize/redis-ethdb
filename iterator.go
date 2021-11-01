package redisethdb

import "github.com/ethereum/go-ethereum/ethdb"

func (db *Database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	panic("implement me")
}
