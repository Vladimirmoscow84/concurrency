/*Условие задачи:
Напишите программу на Go, которая реализует схему обработки данных с использованием паттернов Fan-In и Fan-Out.
Ваше задание:
Реализовать функцию, которая берёт на вход несколько каналов с числами.
Каждые полученные числа нужно параллельно обработать (например, возвести в квадрат или удвоить).
Затем объединим результаты обратно в один канал и выведем их на экран.
Правила:
Использовать паттерн Fan-Out для распределения задач по рабочим горутинам.
Применить паттерн Fan-In для сбора результатов.
Должна быть предусмотрена возможность остановки работы по желанию (через контекст или иными средствами).
Пример:
Пусть у вас есть два канала:

Канал 1: {1, 2, 3}
Канал 2: {4, 5, 6}
Каждый полученный элемент должен быть возведён в квадрат, а результаты собраны в одном канале и напечатаны на экран.

Ожидаемый результат:
Для вышеуказанных каналов результатом будет канал, содержащий числа: {1, 4, 9, 16, 25, 36}.
*/

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func multiply(ctx context.Context, chIn chan int) chan int {
	chOut := make(chan int)
	go func() {
		defer close(chOut)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-chIn:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case chOut <- v * v:
				}
			}
		}
	}()
	return chOut
}
func fanOut(ctx context.Context, chIn chan int, workers int) []chan int {
	channels := make([]chan int, workers)
	for i := range workers {
		channels[i] = multiply(ctx, chIn)
	}

	return channels
}
func fanIn(ctx context.Context, channels ...chan int) chan int {
	resCh := make(chan int)
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
					case resCh <- v:
					}
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(resCh)
	}()
	return resCh
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		ch1 <- 1
		ch1 <- 3
		ch1 <- 5
	}()
	go func() {
		defer close(ch2)
		ch2 <- 2
		ch2 <- 4
		ch2 <- 6
	}()
	const workers = 3
	unitedCh := fanIn(ctx, ch1, ch2)
	channels := fanOut(ctx, unitedCh, workers)
	streamCh := fanIn(ctx, channels...)
	for v := range streamCh {
		fmt.Println(v)
	}

}
