package main

import (
	"fmt"
	"sync"
	"time"
)

type payload struct {
	name  string
	value int
}

func producer(name string, done <-chan struct{}, wg *sync.WaitGroup) <-chan payload {
	ch := make(chan payload)
	i := 1

	go func() {
		defer wg.Done()

		for {
			select {
			case <-done:
				close(ch)
				fmt.Println(name, "completed")
				return
			case ch <- payload{
				name:  name,
				value: i,
			}:
				fmt.Println(name, "produced", i)
				i++
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	return ch
}

func consumer(channels []<-chan payload, done <-chan struct{}, wg *sync.WaitGroup, fanIn chan<- payload) {
	for i, ch := range channels {
		i := i + 1
		ch := ch

		go func() {
			defer wg.Done()
			fmt.Println("Started consume of producer", i)
			for {
				select {
				case <-done:
					fmt.Println("consume of producer", i, "completed")
					return
				case v := <-ch:
					fmt.Println("consumer of producer", i, "got value:", v.value, "from:", v.name)
					fanIn <- v
				}
			}
		}()
	}
}

// применяется когда мы хотим лимитировать параллельные запросы к некоему внешнему сервису
// своего рода рэйт лимитер где в каждый момент времени выполняется только один запрос
func main() {
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	wg.Add(2)
	producers := make([]<-chan payload, 0, 2)
	producers = append(producers, producer("Alice", done, &wg))
	producers = append(producers, producer("Bob", done, &wg))

	fanIn := make(chan payload)

	wg.Add(2)
	consumer(producers, done, &wg, fanIn)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-done:
				return
			case v := <-fanIn:
				fmt.Println("fanIn got:", v)
			}
		}
	}()

	time.Sleep(time.Second)
	close(done)
	wg.Wait()
	time.Sleep(10 * time.Second)
}
