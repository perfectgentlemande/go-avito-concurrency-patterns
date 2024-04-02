package main

import (
	"fmt"
	"sync"
)

const (
	totalJobs    = 10
	totalWorkers = 2
)

func worker(id int, jobs <-chan int, results chan<- int) {
	var wg sync.WaitGroup

	for j := range jobs {
		wg.Add(1)

		go func(job int) {
			defer wg.Done()

			fmt.Println("Worker ", id, " started job ", job)

			result := job * 2
			results <- result

			fmt.Println("Worker ", id, " finished job ", job)
		}(j)
	}

	wg.Wait()
}

// Закон Амдала:
// - в случае когда задача разделяется на несколько частей
// суммарное время ее выполнения на параллельной системе
// не может быть меньше времени выполнения самого медленного фрагмента.
// Если в алгоритме есть последовательная часть обработки,
// то нельзя выполнить алгоритм быстрее
// чем длительность последовательной части алгоритма.

// используется чаще всего для задач выполнения сетевых вызовов

func main() {
	jobs := make(chan int, totalJobs)
	results := make(chan int, totalJobs)

	for w := 1; w <= totalWorkers; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= totalJobs; j++ {
		jobs <- j
	}

	close(jobs)

	for a := 1; a <= totalJobs; a++ {
		<-results
	}
	close(results)
}
