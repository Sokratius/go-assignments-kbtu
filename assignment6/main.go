package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)


// Problem 1


// 1. Sync map
func safeMapSyncMap() {
	fmt.Println("\n--- sync.Map ---")

	var m sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			m.Store("key", val)
		}(i)
	}

	wg.Wait()

	value, _ := m.Load("key")
	fmt.Println("Value:", value)
}

// 2
func safeMapMutex() {
	fmt.Println("\n--- RWMutex ---")

	m := make(map[string]int)
	var mu sync.RWMutex
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			mu.Lock()
			m["key"] = val
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	mu.RLock()
	fmt.Println("Value:", m["key"])
	mu.RUnlock()
}

//Problem 2

// incorrect
func badCounter() {
	fmt.Println("\n--- BAD COUNTER ---")

	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}

// fix 1
func counterMutex() {
	fmt.Println("\n--- COUNTER WITH MUTEX ---")

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}

// fix 2
func counterAtomic() {
	fmt.Println("\n--- COUNTER WITH ATOMIC ---")

	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}

// Problem 3




func startServer(ctx context.Context, name string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
				out <- fmt.Sprintf("[%s] metric: %d", name, rand.Intn(100))
			}
		}
	}()

	return out
}

// FanIn
func FanIn(ctx context.Context, channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	wg.Add(len(channels))

	for _, ch := range channels {
		go func(c <-chan string) {
			defer wg.Done()
			for val := range c {
				select {
				case <-ctx.Done():
					return
				case out <- val:
				}
			}
		}(ch)
	}


	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}


func runFanIn() {
	fmt.Println("\n--- FAN IN ---")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch1 := startServer(ctx, "Alpha")
	ch2 := startServer(ctx, "Beta")
	ch3 := startServer(ctx, "Gamma")

	result := FanIn(ctx, ch1, ch2, ch3)

	for val := range result {
		fmt.Println(val)
	}
}




func main() {

// pr 1
	safeMapSyncMap()
	safeMapMutex()

// pr 2
	badCounter()
	counterMutex()
	counterAtomic()

// pr 3
	runFanIn()

	// The final value is not 1000 because multiple goroutines concurrently update the shared variable without synchronization, causing a race condition.
}