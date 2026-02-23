// fanIn fanOut
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func generator(ctx context.Context, data []int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, value := range data {
			select {
			case <-ctx.Done():
				return
			case out <- value:
			}
		}
	}()

	return out
}

func add(ctx context.Context, chIn <-chan int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-chIn:
				if !ok {
					return
				}
				res := v + 2
				select {
				case <-ctx.Done():
					return
				case <-time.After(200 * time.Millisecond):
					select {
					case <-ctx.Done():
						return
					case out <- res:
					}
				}
			}
		}
	}()
	return out
}
func fanOut(ctx context.Context, n int, chIn <-chan int) []chan int {
	channels := make([]chan int, n)
	for i := range n {
		channels[i] = add(ctx, chIn)
	}

	return channels
}

func fanIn(ctx context.Context, channels []chan int) chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}
	for _, channel := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-channel:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- v:
					}
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	data := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	n := 4 //lenth slice channels
	channels := fanOut(ctx, n, generator(ctx, data))
	streamCh := fanIn(ctx, channels)
	for v := range streamCh {
		fmt.Println("result: ", v)
	}
}
