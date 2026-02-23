// pipeline
package main

import (
	"context"
	"fmt"
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

func add(ctx context.Context, in chan int) chan int {
	out := make(chan int)
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
				result := v + 3
				select {
				case <-ctx.Done():
					return
				case <-time.After(200 * time.Millisecond): //имитация долгой работы
					select {
					case <-ctx.Done():
						return
					case out <- result:

					}

				}
			}
		}
	}()
	return out
}
func mult(ctx context.Context, in chan int) chan int {
	out := make(chan int)
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
				result := v * 2
				select {
				case <-ctx.Done():
					return
				case out <- result:
				}
			}
		}
	}()
	return out
}

func main() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	addCh := add(ctx, generator(ctx, data))
	multCh := mult(ctx, addCh)

	for v := range multCh {
		fmt.Println(v)
	}
	fmt.Println("End of programm")
}
