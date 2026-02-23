/*Напиши реализацию Promise:
Функция promise принимает контекст и функцию, которая возвращает значение и ошибку
Возвращает канал, из которого можно прочитать результат (структура с value и error)
Функция выполняется асинхронно
Контекст должен отменять выполнение*/

package main

import (
	"context"
	"fmt"
	"time"
)

type Result struct {
	Value int
	Err   error
}

func longTask(ctx context.Context) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-time.After(200 * time.Millisecond):
		return 444, nil
	}
}
func promise(ctx context.Context, f func(context.Context) (int, error)) <-chan Result {
	out := make(chan Result, 1)
	go func() {
		defer close(out)
		resCh := make(chan Result)
		go func() {
			value, err := longTask(ctx)
			resCh <- Result{value, err}
		}()

		select {
		case <-ctx.Done():
			out <- Result{0, ctx.Err()}
		case result := <-resCh:
			select {
			case <-ctx.Done():
				out <- Result{0, ctx.Err()}
			case out <- result:
			}
		}

	}()
	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resChan := promise(ctx, longTask)

	result := <-resChan
	if result.Err != nil {
		fmt.Println("Error: ", result.Err)
	} else {
		fmt.Println("result: ", result.Value)
	}
}
