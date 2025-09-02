/*
   Слияние отсортированных каналов
   Описание задачи:

   Даны несколько отсортированных каналов, каждый из которых передает числа в возрастающем порядке. Необходимо объединить
   данные из этих каналов в один канал, который будет содержать все числа также в отсортированном порядке.

   Условия:
   - Каналы, поданные на вход функции, передают числа в возрастающем порядке.
   - Необходимо вернуть новый канал, который будет содержать все числа из этих каналов, также в отсортированном порядке.
   - Обработка данных должна быть конкурентной.
*/

package main

import (
	"fmt"
)

//струууутура для хранения элементов каналов

type Item struct {
	value int
	ch    chan int
}

// Определяем тип для реализации минимальной кучи
type MiniHeap []Item

// mergeSortedChannels объединяет несколько отсортированных каналов в один отсортированный канал
func mergeSortedChannels(channels ...chan int) chan int {
	// your code here
	resCh := make(chan int)

	return resCh
}

func main() {
	// Пример использования
	chan1 := make(chan int)
	chan2 := make(chan int)
	chan3 := make(chan int)

	// Эмуляция передачи данных в каналы
	go func() {
		defer close(chan1)
		chan1 <- 1
		chan1 <- 3
		chan1 <- 5
	}()

	go func() {
		defer close(chan2)
		chan2 <- 2
		chan2 <- 4
	}()

	go func() {
		defer close(chan3)
		chan3 <- 6
		chan3 <- 8
	}()

	// Вызываем функцию mergeSortedChannels и выводим результат
	mergedChan := mergeSortedChannels(chan1, chan2, chan3)
	for num := range mergedChan {
		fmt.Println(num) // Ожидаемый вывод: 1, 2, 3, 4, 5, 6, 8
	}
}
