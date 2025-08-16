package store

import "sync"

type Memory struct {
	wg sync.Mutex
}

type IMemory interface {
	Get(dbName string)
	Set(T any, db string)
	Delete(db string, key string)
}
