/*
Задача повышенной сложности на семафор в Go
Условие задачи:
Вы разрабатываете высоконагруженный сервис резервирования авиабилетов, работающий в режиме реального времени. Сервер получает огромное количество запросов на бронирование одновременно, и необходимо обеспечить корректную обработку запросов с соблюдением определенных ограничений:
-Система может одновременно обрабатывать не более 100 запросов на бронирование.
-Каждый запрос обрабатывается индивидуально, и процесс бронирования длится от 1 до 3 секунд.
-При превышении лимита одновременно обрабатываемых запросов лишние запросы становятся в очередь и ожидают освобождения слота.

Важно поддерживать статистику: количество выполненных запросов, среднее время обработки и процент успешности обработки.
Задачи:
--Реализуйте службу бронирования с использованием семафора, который ограничивает количество одновременно обрабатываемых запросов.
--Обеспечьте ведение статистики:
--Сколько запросов выполнено.
--Среднее время обработки запроса.
--Процент успешных запросов (предполагая, что 5% запросов имеют ошибку обработки).
--Разработайте программу, которая демонстрирует стабильную работу службы при высоком уровне нагрузки.

Ограничения:
---Сервер должен выдерживать нагрузку до 1000 запросов в минуту.
---Запросы должны обрабатываться в порядке поступления.
---Семафор должен использоваться корректно, чтобы предотвратить перегрузку сервера.
*/
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type order struct {
	orderId    int
	clientId   int
	durationMs int
	err        error
}

var mu sync.Mutex

func reserveTickets(semaphore chan struct{}, reqCh chan order, clientId int, orderCount *int, wg *sync.WaitGroup) {
	defer wg.Done()
	semaphore <- struct{}{}
	now := time.Now()
	time.Sleep(time.Duration(rand.Intn(100) * int(time.Millisecond)))
	mu.Lock()
	*orderCount++
	orderId := *orderCount
	mu.Unlock()

	if rand.Float64() < 0.05 {
		reqCh <- order{
			orderId:    orderId,
			clientId:   clientId,
			durationMs: int(time.Since(now).Milliseconds()),
			err:        errors.New("error of request"),
		}
		fmt.Printf("error of request in order %d from client %d\n", orderId, clientId)
		<-semaphore
		return
	}

	reqCh <- order{
		orderId:    orderId,
		clientId:   clientId,
		durationMs: int(time.Since(now).Milliseconds()),
		err:        nil,
	}

	fmt.Printf("successful request in order %d from client %d\n", orderId, clientId)
	<-semaphore
}
func main() {
	orderCount := new(int)
	clientCount := 200
	reqCh := make(chan order)

	reqNum := 3
	semaphore := make(chan struct{}, reqNum)
	wg := &sync.WaitGroup{}

	for i := 1; i <= clientCount; i++ {
		wg.Add(1)
		go reserveTickets(semaphore, reqCh, i, orderCount, wg)
	}

	go func() {
		wg.Wait()
		close(reqCh)
	}()

	errOrderCount := 0
	dur := 0
	for ord := range reqCh {
		if ord.err != nil {
			errOrderCount++
		}
		dur = dur + ord.durationMs
	}

	percent := errOrderCount * 100 / (*orderCount)
	dursr := dur / (*orderCount)

	fmt.Printf("percent errors: %d\n", percent)
	fmt.Printf("dur: %d\n", dursr)
}
