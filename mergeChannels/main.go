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
	"sync"
	"time"
)

func mergeChannels(ctx context.Context, channels ...chan string) chan string {
	chOut := make(chan string)

	wg:=sync.WaitGroup{}
	mu:=sync.Mutex{}

	for _,channel:=range channels{
		wg.Go(func(){
			for{
				select{
				case<-ctx.Done():
					return
				case v,ok:=<-channel:
					if !ok{
						return
					}
					select{
					case <-ctx.Done():
						return
						case 
					}
				}
			}
		})
	}

	return chOut
}
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	go func() {
		defer close(ch1)
		ch1 <- "ssd"
		ch1 <- "qwe"
		ch1 <- "spartak"
		ch1 <- "ssd"
	}()

	go func() {
		defer close(ch2)
		ch2 <- "opi"
		ch2 <- "jgjhghlklh"
		ch2 <- "spartak1"
		ch2 <- "qwe"
	}()
	go func() {
		defer close(ch3)
		ch3 <- "qwerty"
		ch3 <- "dfghgfd"
		ch3 <- "spartak2"
		ch3 <- "zalupa"
	}()

}
