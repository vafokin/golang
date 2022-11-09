package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	storage map[string]int
	mu      sync.RWMutex
}

const (
	k1   = "key1"
	step = 7
)

func (c *Cache) Increase(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage[key] += value
}

func (c *Cache) Set(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage[key] = value
}

func (c *Cache) Get(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.storage[key]
}

func (c *Cache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.storage, key)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	semaphore := make(chan int, 4)
	cache := Cache{storage: make(map[string]int)}
	defer cancel()

	for i := 0; i < 10; i++ {
		select {
		case semaphore <- i:
			go func() {
				select {
				case <-ctx.Done():
					fmt.Println("context deadline exceeded")
					return
				default:
					defer func() {
						msg := <-semaphore
						fmt.Println(msg)
					}()
					cache.Increase(k1, step)
					time.Sleep(time.Millisecond * 100)
				}
			}()
			/*for len(semaphore) > 0 {
				time.Sleep(time.Millisecond * 10)
			}*/
		case <-ctx.Done():
			fmt.Println("context deadline exceeded1")
			return
		}
	}

	for i := 0; i < 10; i++ {
		select {
		case semaphore <- i:
			go func(i int) {
				select {
				case <-ctx.Done():
					fmt.Println("context deadline exceeded2")
					return
				default:
					defer func() {
						msg := <-semaphore
						fmt.Println(msg)
					}()
					cache.Set(k1, step*i)
					time.Sleep(time.Millisecond * 100)
				}
			}(i)
		case <-ctx.Done():
			fmt.Println("context deadline exceeded3")
			return
			/*for len(semaphore) > 0 {
				time.Sleep(time.Millisecond * 10)
			} */
		}
	}
	fmt.Println(cache.Get(k1))
}
