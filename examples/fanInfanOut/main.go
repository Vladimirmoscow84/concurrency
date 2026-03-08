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
		for _, v := range data {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return out
}
func fanOut(ctx context.Context, n int, in chan int) []chan int {
	channels := make([]chan int, n)
	for i := 0; i < n; i++ {
		out := make(chan int)
		channels[i] = out
		go func() {
			defer close(out)
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					res := v * 2
					select {
					case <-ctx.Done():
						return
					case <-time.After(120 * time.Millisecond):
						select {
						case <-ctx.Done():
							return
						case out <- res:
						}
					}
				}
			}
		}()
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
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	genCh := generator(ctx, data)
	n := 4 //number workers
	channels := fanOut(ctx, n, genCh)
	streamCh := fanIn(ctx, channels)

	for v := range streamCh {
		fmt.Println(v)
	}

	fmt.Println("End of Programm")

}
