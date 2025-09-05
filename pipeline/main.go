/*  патерн, который разделяет сложную задачу на более легкие, которые выполняются в отдельных горутинах*/
package main

import (
	"context"
	"fmt"
	"time"
)

func generator(ctx context.Context, data []int) chan int {
	chOut := make(chan int)
	go func() {
		defer close(chOut)
		for _, v := range data {
			select {
			case <-ctx.Done():
				fmt.Println("done in generator")
				return
			case chOut <- v:
			}
		}
	}()

	return chOut
}
func add(ctx context.Context, chIn chan int) chan int {
	chOut := make(chan int)
	go func() {
		defer close(chOut)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("done in add function")
				return
			case v, ok := <-chIn:
				if !ok {
					return
				}
				result := v + 2
				select {
				case <-ctx.Done():
					fmt.Println("done in add function")
					return
				case chOut <- result:

				}

			}
		}
	}()
	return chOut
}

func multiply(ctx context.Context, chIn chan int) chan int {
	chOut := make(chan int)
	go func() {
		defer close(chOut)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("done in multiply function")
				return
			case v, ok := <-chIn:
				if !ok {
					return
				}
				result := v * v
				select {
				case <-ctx.Done():
					return
				case chOut <- result:
				}

			}
		}
	}()
	return chOut
}
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Second))
	defer cancel()
	data := []int{1, 2, 3, 4, 5, 6, 7}
	chData := generator(ctx, data)
	resCh := multiply(ctx, add(ctx, chData))

	for v := range resCh {
		fmt.Println(v)
		time.Sleep(300 * time.Millisecond)
	}
}
