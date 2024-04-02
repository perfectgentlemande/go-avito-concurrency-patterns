package main

import (
	"fmt"
)

type ringBuffer struct {
	inCh  chan int
	outCh chan int
}

func NewRingBuffer(inCh, outCh chan int) *ringBuffer {
	return &ringBuffer{
		inCh:  inCh,
		outCh: outCh,
	}
}

func (r *ringBuffer) Run() {
	defer close(r.outCh)

	for v := range r.inCh {
		select {
		case r.outCh <- v:
		default:
			<-r.outCh
			r.outCh <- v
		}
	}
}

// используется в основном в случаях
// когда нужно успевать обрабатывать входной поток данных
// где некоторые элементы входного потока могут быть без последствия отброшены
// таким образом прореживая слишком быстрый поток
// если успеваем все обработать, то прореживания не происходит
// например при отправке большого количества метрик чтобы не нагружать хранилище метрик

func main() {
	inCh := make(chan int)
	outCh := make(chan int, 4)

	rb := NewRingBuffer(inCh, outCh)
	go rb.Run()

	for i := 1; i <= 10; i++ {
		inCh <- i
	}

	close(inCh)

	// без этой паузы считаем больше буффера данных, а с ней только по размеру буффера
	// time.Sleep(100 * time.Millisecond)

	for res := range outCh {
		fmt.Println(res)
	}
}
