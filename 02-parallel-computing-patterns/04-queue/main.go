package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	N        = 3
	MESSAGES = 10
)

func process(payload int, queue chan struct{}, wg *sync.WaitGroup) {
	queue <- struct{}{}

	go func() {
		defer wg.Done()

		fmt.Println("Start processing of ", payload)
		time.Sleep(time.Millisecond * 500)
		fmt.Println("Completed processing of ", payload)

		<-queue
	}()
}

// позволяет принимать до n сообщений не ожидая их обработки
// везде где есть отложенные задачи выполнения которых мы не ожидаем
// например отправка почты

func main() {
	var wg sync.WaitGroup

	fmt.Println("Queue of length N:", N)
	queue := make(chan struct{}, N)

	wg.Add(MESSAGES)

	for w := 1; w <= MESSAGES; w++ {
		process(w, queue, &wg)
	}

	wg.Wait()

	close(queue)
	fmt.Println("Process completed")
}
