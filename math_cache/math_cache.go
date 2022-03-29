// A caching method for challenge.
//
// Code example acquired from https://hackernoon.com/in-memory-caching-in-golang and adjusted for the challenge

package math_cache

import (
	"fmt"
	"sync"
	"time"
)

type cachedOperation struct {
	answer            int
	expireAtTimestamp int64
}

func operation_id(operation string, x int, y int) string {
	return fmt.Sprint(operation, x, y)
}

type localCache struct {
	stop       chan struct{}
	wg         sync.WaitGroup
	mu         sync.RWMutex
	operations map[string]cachedOperation
}

var (
	lc *localCache
)

func NewLocalCache() *localCache {

	if lc == nil {
		lc = &localCache{
			operations: make(map[string]cachedOperation),
		}

		lc.wg.Add(1)
		go func(cleanupInterval time.Duration) {
			defer lc.wg.Done()
			lc.cleanupLoop(cleanupInterval)
		}(time.Minute * 2)
	}
	return lc
}

func (lc *localCache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	fmt.Println("Cache cleanup started")
	for {
		select {
		case <-lc.stop:
			return
		case <-t.C:
			lc.mu.Lock()
			for op_id, co := range lc.operations {
				if co.expireAtTimestamp <= time.Now().Unix() {
					fmt.Println(fmt.Sprintf("Removed id '%s'", op_id))
					delete(lc.operations, op_id)
				}
			}
			lc.mu.Unlock()
		}
	}
}

func (lc *localCache) extendExpireAt() int64 {
	return time.Now().Add(time.Minute * 1).Unix()
}

func (lc *localCache) Update(operation string, x int, y int, answer int) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	op_id := operation_id(operation, x, y)
	lc.operations[op_id] = cachedOperation{
		answer:            answer,
		expireAtTimestamp: lc.extendExpireAt(),
	}
}

func (lc *localCache) Read(operation string, x int, y int) (int, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	op_id := operation_id(operation, x, y)

	co, ok := lc.operations[op_id]
	if !ok {
		fmt.Println(fmt.Sprintf("Math operation '%s' with operands %d and %d is not in cache.", operation, x, y))
		return 0, false
	}

	if co.expireAtTimestamp <= time.Now().Unix() {
		fmt.Println(fmt.Sprintf("Math operation '%s' with operands %d and %d cache is expired.", operation, x, y))
		delete(lc.operations, operation_id(operation, x, y))
		return 0, false
	}
	fmt.Println(fmt.Sprintf("Math operation '%s' with operands %d and %d cache hit.", operation, x, y))
	return co.answer, true
}
