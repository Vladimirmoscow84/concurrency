// generator
package main

import (
	"context"
	"fmt"
	"time"
)

func generator(ctx context.Context, data []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, v := range data {
			select {
			case <-ctx.Done():
				fmt.Printf("exit with context %v", ctx.Err())
				return
			case out <- v:
			}
		}

	}()

	return out
}
func main() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for v := range generator(ctx, data) {
		fmt.Println(v)
	}
	fmt.Println("End of programm")
}
