package main

import (
	"context"
	"fmt"
	"sync"

	"time"
)

func worker(ctx context.Context, n int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Ewit wit context : %v\n", ctx.Err())
			return
		case v, ok := <-jobs:
			if !ok {
				//fmt.Println("Channel jobs is empty")
				return
			}
			fmt.Printf("worker %d begin job %d\n", n, v)
			res := v + 2
			time.Sleep(300 * time.Millisecond)
			select {
			case <-ctx.Done():
				fmt.Printf("Ewit wit context : %v\n", ctx.Err())
				return
			case results <- res:
				fmt.Printf("worker %d finish job %d\n", n, v)
			}
		}

	}

}
func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	numJ := 20
	numW := 5
	jobs := make(chan int, numJ)
	results := make(chan int, numJ)

	for i := 1; i <= numW; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for j := 1; j <= numJ; j++ {
			select {
			case <-ctx.Done():
				fmt.Println("exit with context", ctx.Err())
				close(jobs)
				return
			case jobs <- j:
				fmt.Printf("job %d has been sended\n", j)
			}
		}
		close(jobs)
	}()

	for res := range results {
		fmt.Println("result is: ", res)
	}

}
