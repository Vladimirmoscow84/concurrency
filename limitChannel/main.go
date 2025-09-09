/*
Ограничитель скорости для канала
Описание задачи:

Дана функция-генератор, которая непрерывно пишет числа в канал.
Необходимо реализовать функцию limitChannel, которая принимает этот канал
и возвращает новый канал, но с ограничением скорости передачи данных.

Условия:
- Функция limitChannel принимает входной канал и максимальное количество сообщений в секунду (rate).
- Должен вернуться новый канал, в который попадают данные из исходного, но не чаще, чем rate раз в секунду.
- Если исходный канал закрывается — новый канал также должен закрыться.
- Обработка должна выполняться конкурентно.
*/
package main

import (
	"context"
	"fmt"
	"time"
)

func generator(ctx context.Context) chan int {
	chOut := make(chan int)
	go func() {
		defer close(chOut)
		v := 0
		for {
			select {
			case <-ctx.Done():
				return
			case chOut <- v:
				v++
			}

		}

	}()
	return chOut
}
func limmitChannel(ctx context.Context, inputCh chan int, rate int) chan int {
	chStream := make(chan int)

	interval := time.Second / time.Duration(rate)
	ticker := time.NewTicker(interval)
	go func() {
		defer close(chStream)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-inputCh:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					select {
					case <-ctx.Done():
						return
					case chStream <- v:
					}
				}
			}

		}

	}()

	return chStream
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	chGen := generator(ctx)
	finalStream := limmitChannel(ctx, chGen, 5)

	for v := range finalStream {
		fmt.Println(v)
	}

}
