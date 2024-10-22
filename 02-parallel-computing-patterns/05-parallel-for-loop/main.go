package main

import (
	"fmt"
	"time"
)

type empty struct{}

const DATA_SIZE = 4

func calculate(val int) int {
	time.Sleep(time.Millisecond * 500)
	return val * 2
}

// когда получаем бэтч данных и с ним сходить куда-то еще, данные в батчи не зависят друг от друга
func main() {
	data := make([]int, 0, DATA_SIZE)
	for i := 0; i < DATA_SIZE; i++ {
		data = append(data, i+10)
	}
	results := make([]int, DATA_SIZE)
	semaphore := make(chan empty, DATA_SIZE)

	fmt.Println("Before: ", data)
	start := time.Now()

	for i, xi := range data {
		go func(i int, xi int) {
			results[i] = calculate(xi)
			semaphore <- empty{}
		}(i, xi)
	}

	for i := 0; i < DATA_SIZE; i++ {
		<-semaphore
	}

	fmt.Println("After: ", results)
	fmt.Println("Elapsed: ", time.Since(start))
}
