/*
Простая задача на пулл воркеров (Pool Workers) в Go
Условие задачи:
Вы разрабатываете систему, которая обрабатывает задачи в отдельном пуле воркеров (рабочих потоков). Предположим, что у вас есть некоторый набор задач, которые нужно выполнить параллельно. Вам необходимо создать пулл воркеров, состоящий из фиксированного количества горутин, каждая из которых способна обрабатывать поступающие задачи.

Требования:
Создайте пулл воркеров, состоящий из 3 рабочих потоков (воркеров).
Пулл должен быть способен обрабатывать задачи из входящего канала.
Воркер выполняет простую задачу: получает число, удваивает его и отправляет результат в выходной канал.
Запустите пулл и продемонстрируйте его работу на примерах.
Пример работы программы:
На вход подается канал с числами: {.....}.
Программа должна обработать эти числа и возвратить удвоенные значения: {.....}.
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

// poolOfWorkers реализует пулл воркеров, обрабатывающих задачи.
func poolOfWorkers(tasks chan int, results chan int, workers int) {
	wg := sync.WaitGroup{}
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range tasks {

				result := v * 2
				results <- result
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
}

func main() {
	now := time.Now()
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	tasks := make(chan int)
	results := make(chan int)
	workers := 3

	poolOfWorkers(tasks, results, workers)

	go func() {
		defer close(tasks)
		for _, v := range data {
			tasks <- v
		}
	}()

	for answer := range results {
		fmt.Println(answer)
	}

	fmt.Printf("Работа программы завершена, время работы: %v", time.Since(now))
}
