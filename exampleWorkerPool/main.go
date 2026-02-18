package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, jobs <-chan int, results chan<- int, w int) {
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Printf("worker %d begin job %d\n", w, v)
			time.Sleep(1 * time.Second)
			result := v * 2
			select {
			case <-ctx.Done():
				return
			case results <- result:
				fmt.Printf("worker %d finished job %d\n", w, v)
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	numJ := 20 //num jobs
	numW := 5  //num workers
	jobs := make(chan int, numJ)
	results := make(chan int, numJ)

	for i := 1; i <= numW; i++ {
		go worker(ctx, jobs, results, i)
	}

	for j := 1; j <= numJ; j++ {
		select {
		case <-ctx.Done():
			return
		case jobs <- j:

		}
	}
	close(jobs)

	for result := range results {
		fmt.Println("result: ", result)
	}
}
