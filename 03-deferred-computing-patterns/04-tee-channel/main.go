package main

import (
	"context"
	"fmt"
	"time"
)

func tee(ctx context.Context, input <-chan string, outputs []chan<- string) {
	for elem := range input {
		elem := elem

		for _, out := range outputs {
			out := out

			go func() {
				select {
				case out <- elem:
					break
				case <-ctx.Done():
					break
				}
			}()
		}
	}
}

// вроде и похож на fan-out однако сообщения дублируются во все каналы
// можно использовать при реализации брокера сообщения
// чтобы его продублировать в несколько потребителей
func main() {
	in := make(chan string)
	out1 := make(chan string)
	out2 := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())

	outputs := []chan<- string{out1, out2}

	go func() {
		in <- "A"
		in <- "B"
		in <- "C"
		close(in)
	}()

	tee(ctx, in, outputs)
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case v := <-out1:
				fmt.Println("out1 got value: ", v)
			case v := <-out2:
				fmt.Println("out2 got value: ", v)
			}
		}
	}()

	time.Sleep(time.Second)
	cancel()
}
