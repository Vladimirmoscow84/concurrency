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
	"container/heap"
	"fmt"
)

//струууутура для хранения элементов каналов

type Item struct {
	value int
	ch    chan int
}

// Определяем тип для реализации минимальной кучи
type MiniHeap []Item

// опрепделяем методы лоя кучи

// Len - возвращает длину
func (h MiniHeap) Len() int {
	return len(h)
}

// Less - определяет, какой элемент меньше другого
func (h MiniHeap) Less(i, j int) bool {
	return h[i].value < h[j].value
}

// Swap - перестановка элементов местами по указанным индексам
func (h MiniHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push - добавляет эелемент в кучу
func (h *MiniHeap) Push(x any) {
	*h = append(*h, x.(Item))
}

// Pop - удаляет и возвращает последний элемент из кучи
func (h *MiniHeap) Pop() any {
	// oldHeap := *h
	// n := len(oldHeap)
	// item := oldHeap[n-1]
	// *h = oldHeap[:n-1]
	// return item

	item := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return item
}

// mergeSortedChannels объединяет несколько отсортированных каналов в один отсортированный канал
func mergeSortedChannels(channels ...chan int) chan int {
	// your code here
	out := make(chan int)
	go func() {
		defer close(out)

		// Создание кучи для первого элемента каждого канала
		h := &MiniHeap{}
		heap.Init(h)
		for _, ch := range channels {
			v, ok := <-ch
			if !ok {
				continue // Пропустить закрытые каналы
			}
			heap.Push(h, Item{v, ch})
		}

		// Начинаем сливать каналы
		for h.Len() > 0 {
			// Берём минимальный элемент из кучи
			top := heap.Pop(h).(Item)
			out <- top.value

			// Пробуем прочитать следующий элемент из канала
			nextVal, ok := <-top.ch
			if ok {
				// Если есть следующее значение, добавляем его в кучу
				heap.Push(h, Item{nextVal, top.ch})
			}
		}
	}()
	return out
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
