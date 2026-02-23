/*
Напиши реализацию Promise:
Функция promise принимает контекст и функцию, которая возвращает значение и ошибку
Возвращает канал, из которого можно прочитать результат (структура с value и error)
Функция выполняется асинхронно
Контекст должен отменять выполнение
*/
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
		fRes := make(chan Result)
		go func() {
			val, err := f(ctx)
			select {
			case <-ctx.Done():
				return
			case fRes <- Result{val, err}:
			}
		}()

		select {
		case <-ctx.Done():
			out <- Result{0, ctx.Err()}
		case value := <-fRes:
			select {
			case <-ctx.Done():
				out <- Result{0, ctx.Err()}
			case out <- value:
			}
		}

	}()
	return out
}
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resCh := promise(ctx, longTask)

	result := <-resCh
	if result.Err != nil {
		fmt.Println("Error: ", result.Err)
	} else {
		fmt.Println("result: ", result.Value)
	}
	fmt.Println("End of programm")
}
