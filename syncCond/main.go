// написать обычный пример использования примитива sync.Cond
package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	cond := sync.NewCond(&mu)
	condition := false
	ready := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		ready <- struct{}{}
		fmt.Println("1 - Горутина 1 запустилась")
		mu.Lock()
		defer mu.Unlock()
		for !condition {
			fmt.Println("3 -горутина 1 ждет сигнал из мэйн, освобождая мютекс и засыпает")
			cond.Wait()
		}
		fmt.Println("5 - горутина 1 получила сигнал из мэйн и проснулась")

	}()
	<-ready
	fmt.Println("2 - горутина 1 точно запустилась")
	mu.Lock()
	condition = true
	fmt.Println("4 - сигнал из мэйн отправлен в горутину 1")
	cond.Signal()
	mu.Unlock()

	wg.Wait()
	fmt.Println("END!!!")
}
