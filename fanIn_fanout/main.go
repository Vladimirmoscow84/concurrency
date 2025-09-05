/*
Паттерны отвечающие за:
fan-out - принимает канал с входными данными и запускает заданное множество горутин для одновременной обработки этих данных
fan-in - обединяет несколько результатов c с разных канало в один канал (мультиплесирование)
*/
package main

import (
	"context"
	"fmt"
	"time"
)

// genrator - читает входные данные и отправляет их в канал
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
	resChan := make(chan int)
	go func() {
		defer close(resChan)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("done in add function")
				return
			case v, ok := <-chIn:
				if !ok {
					return
				}
				result := v + 13
				select {
				case <-ctx.Done():
					fmt.Println("done in add function")
				case resChan <- result:
				}
			}
		}

	}()

	return resChan
}
func main() {
	//задано какое-то множество данных в виде слайса
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Second))
	defer cancel()
}
