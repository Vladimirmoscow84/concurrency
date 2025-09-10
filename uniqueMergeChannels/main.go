/*
Вам даны несколько каналов, каждый из которых отправляет случайные числа в диапазоне от 1 до 100. Ваша задача — написать функцию, которая объединяет все эти каналы в один канал таким образом, чтобы в итоге все числа из входных каналов попадали в единый выходной канал без дублирования и сохраняли уникальный порядок следования (первое появление числа считается уникальным порядком).

Пример:
Предположим, имеются три канала:

Канал 1: {1, 3, 5}
Канал 2: {2, 4, 5}
Канал 3: {3, 6, 7}
Результатом вашей функции будет канал, возвращающий следующую последовательность чисел: {1, 2, 3, 4, 5, 6, 7}, без повторений.
*/
package main

import (
	"context"
	"sync"
	"time"
)

// для сортировки поступающих данных из каналов создаем миниКучу, состоящую из структур типа(данные-канал)
// для проверки повторяемости данных создаем хранилище (мапу) уникальных значений данных из каналов
func uniqueMergeChannels(ctx context.Context, channels ...chan int) chan int {
	chOut := make(chan int)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	storage := make(map[int]struct{})

	return chOut
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)

	go func() {
		defer close(ch1)
		ch1 <- 3
		ch1 <- 22
		ch1 <- 23
		ch1 <- 30
	}()
	go func() {
		defer close(ch2)
		ch2 <- 2
		ch2 <- 3
		ch2 <- 4
		ch2 <- 7
		ch2 <- 23
	}()
	go func() {
		defer close(ch3)
		ch3 <- 4
		ch3 <- 7
		ch3 <- 12
		ch3 <- 22
	}()
	go func() {
		defer close(ch4)
		ch4 <- 2
		ch4 <- 3
		ch4 <- 7
		ch4 <- 15
		ch4 <- 16
		ch4 <- 17
		ch4 <- 23
		ch4 <- 25
	}()

}
