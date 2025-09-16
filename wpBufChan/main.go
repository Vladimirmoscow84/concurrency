/*
Задача повышенной сложности на пулл воркеров (Pool Workers) в Go
Условие задачи:
Реализуйте систему обработки запросов на расчёт площади прямоугольника. При этом вы должны учесть следующие факторы:

Количество воркеров: Пулл воркеров должен динамически увеличиваться или уменьшаться в зависимости от текущей нагрузки.
Буферизация: Входящий канал запросов должен быть буферизированным, чтобы временно накапливать запросы, если воркеры заняты.
Таймаут: Если запрос не обработался в течение определённого времени (например, 5 секунд), он должен считаться просроченным и выводиться отдельное предупреждение.
Регистрация успешной обработки: После обработки запроса программа должна выводить подробную информацию о нём.
Дополнительные требования:
При достижении порогового уровня нагрузки (например, количество незавершённых запросов достигает 10) количество воркеров увеличивается на единицу.
Если нагрузка снижается (все запросы обработались), количество воркеров уменьшается до оптимального значения (например, до 3).
Каждый запрос должен содержать поля: ID запроса, ширину и высоту прямоугольника.
Пример работы программы:
Исходные данные: запросы на расчет площадей прямоугольников, представленные в виде структуры {id, width, height}.

	requests := []struct{id int; width, height float64}{
	    {1, 10, 20},
	    {2, 15, 15},
	    ...
	}

Желаемый результат:

Request #1 processed: area = 200
Request #2 processed: area = 225
...
Warning: Request #5 timed out after 5 seconds.
All requests processed or timed out.
*/
package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Request struct {
	id   int
	x, y float64
}
type Result struct {
	request Request
	area    float64
	err     error
}

type PoolManager struct {
	reqCh   chan Request
	resCh   chan Result
	doneCh  chan struct{}
	workers uint32
	wg      sync.WaitGroup
}

const (
	workers     = 3
	buffSize    = 10
	loadService = 10
	timeout     = 5 * time.Second
)

func newPoolManager() *PoolManager {

	return &PoolManager{
		reqCh:   make(chan Request, buffSize),
		resCh:   make(chan Result),
		doneCh:  make(chan struct{}),
		workers: workers,
	}
}
func (pm *PoolManager) worker() {
	defer pm.wg.Done()
	for {
		select {
		case req, ok := <-pm.reqCh:
			if !ok {
				return
			}
			timer := time.AfterFunc(timeout, func() {
				pm.resCh <- Result{req, 0, fmt.Errorf("time is over")}
			})
			area := req.x * req.y
			timer.Stop()
			pm.resCh <- Result{req, area, nil}

		case <-pm.doneCh:
			return

		}
	}
}
func (pm *PoolManager) start() {
	for i := 1; i <= int(atomic.LoadUint32(&pm.workers)); i++ {
		pm.wg.Add(1)
		go pm.worker()
	}

	log.Println("Pool manager started")
}

func (pm *PoolManager) stop() {
	close(pm.doneCh)
	pm.wg.Wait()
	log.Println("Pool manage stopped")
}
func (pm *PoolManager) getresults() chan Result {

	return pm.resCh
}
func main() {

}
