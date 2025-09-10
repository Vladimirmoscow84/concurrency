/*
   Слияние каналов строк с удалением дубликатов
   Описание задачи:

   Даны несколько каналов строк, каждый из которых передает строки в произвольном порядке.
   Необходимо объединить данные из этих каналов в один канал, но при этом каждая строка
   должна встречаться в выходном канале только один раз.

   Условия:
   - Функция должна принимать на вход несколько каналов строк.
   - Должен вернуться новый канал, из которого можно читать объединённые строки без повторов.
   - Порядок появления строк в результирующем канале не фиксируется (может быть любым).
   - Обработка должна быть конкурентной.
*/

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func mergeChannels(ctx context.Context, channels ...chan string) chan string {
	chOut := make(chan string)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	storage := make(map[string]bool)

	for _, channel := range channels {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("exit with context Timeout")
					return
				case v, ok := <-channel:
					if !ok {
						return
					}
					mu.Lock()
					if _, exist := storage[v]; !exist {
						storage[v] = true
						select {
						case <-ctx.Done():
							fmt.Println("exit with context Timeout")
							return
						case chOut <- v:
						}
					}
					mu.Unlock()
				}
			}
		})

	}
	go func() {
		wg.Wait()
		close(chOut)
	}()

	return chOut
}
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	go func() {
		defer close(ch1)
		ch1 <- "ssd"
		//time.Sleep(500 * time.Millisecond)
		ch1 <- "qwe"
		//time.Sleep(500 * time.Millisecond)
		ch1 <- "spartak"
		//time.Sleep(500 * time.Millisecond)
		ch1 <- "ssd"
	}()

	go func() {
		defer close(ch2)
		ch2 <- "opi"
		ch2 <- "jgjhghlklh"
		//time.Sleep(200 * time.Millisecond)
		ch2 <- "spartak1"
		ch2 <- "qwe"
		//time.Sleep(500 * time.Millisecond)
	}()
	go func() {
		defer close(ch3)
		//time.Sleep(200 * time.Millisecond)
		ch3 <- "qwerty"
		ch3 <- "dfghgfd"
		//time.Sleep(700 * time.Millisecond)
		ch3 <- "spartak2"
		//time.Sleep(200 * time.Millisecond)
		ch3 <- "zalupa"
	}()

	stream := mergeChannels(ctx, ch1, ch2, ch3)

	for v := range stream {
		fmt.Println(v)
	}

}
