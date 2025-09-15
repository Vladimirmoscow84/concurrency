/*
Задача на пулл воркеров в Go
Условие задачи:
Разработайте программу, которая рассчитывает произведение двух чисел. Допустим, у вас есть список пар чисел, и вам нужно параллельно рассчитать произведения каждой пары, используя пулл воркеров.

Требования:

Количество воркеров фиксированное (например, 3).
Каждый воркер получает задачу из входящего канала, обрабатывает её и отправляет результат в выходной канал.
Необходимо правильно организовать обработку данных и закрыть каналы по завершении работы.
Пример работы программы:
Исходные данные: пара чисел {1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}.

Желаемый результат:

Product of 1 and 2 is 2
Product of 3 and 4 is 12
Product of 5 and 6 is 30
Product of 7 and 8 is 56
Product of 9 and 10 is 90
All tasks completed.
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Task struct {
	x, y, z int
}

func workerPool(tasks chan Task, results chan Task, workers int) {
	wg := sync.WaitGroup{}

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				task.z = task.x * task.y
				results <- task
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
}

func main() {
	data := make([]Task, 0)
	for range 500 {
		data = append(data, Task{rand.Intn(3000), rand.Intn(3000), 0})
	}
	workers := 3
	tasks := make(chan Task)
	results := make(chan Task)

	workerPool(tasks, results, workers)

	go func() {
		for _, value := range data {
			tasks <- value
		}
		close(tasks)
	}()

	for v := range results {
		fmt.Printf("Product of %d and %d = %d\n", v.x, v.y, v.z)
	}
}
