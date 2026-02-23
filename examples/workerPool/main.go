/*Напиши паттерн Worker Pool. Условия:
Функция worker принимает контекст, задания из канала, результаты в канал
В main запускается N воркеров
M заданий отправляется в канал
Результаты собираются и выводятся
Контекст с таймаутом 2 секунды*/

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
			return
		case v, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Printf("worker %d begin job %d\n", n, v)
			res := v * 2
			select {
			case <-ctx.Done():
				return
			case <-time.After(200 * time.Millisecond):
				fmt.Printf("worker %d finish job %d\n", n, v)
				select {
				case <-ctx.Done():
					return
				case results <- res:
				}

			}
		}
	}
}
func main() {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	n := 5  //num workers
	m := 20 //num jobs
	jobs := make(chan int, m)
	results := make(chan int, m)

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for j := 1; j <= m; j++ {
		jobs <- j
	}
	close(jobs)

	for res := range results {
		fmt.Println("result: ", res)
	}
	fmt.Println("End of programm")
}
