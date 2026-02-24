//semaphore

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Sema struct {
	sema chan struct{}
}

func newSema(cap int) *Sema {
	if cap <= 0 {
		panic("capacity must be positive")
	}
	return &Sema{
		sema: make(chan struct{}, cap),
	}
}
func (s *Sema) AcquireWithContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.sema <- struct{}{}:
		return nil
	}
}

func (s *Sema) Release() {
	<-s.sema
}

func main() {
	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c := 10  //sema capacity
	w := 100 //num workers
	wg := sync.WaitGroup{}
	Sema1 := newSema(c)

	for i := 1; i <= w; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := Sema1.AcquireWithContext(ctx)
			if err != nil {
				fmt.Printf("worker have not token %v\n", err)
				return
			}
			defer Sema1.Release()
			fmt.Printf("worker %d begin job\n", i)
			result := 2 * i
			time.Sleep(150 * time.Millisecond)
			fmt.Printf("worker %d finish job, result: %d\n", i, result)

		}()
	}
	wg.Wait()
	fmt.Println("end of work", time.Since(now))
}
