package main

import (
	"fmt"
	"time"
)

type data struct {
	Body  string
	Error error
}

func doGet(url string) (string, error) {
	time.Sleep(time.Microsecond * 200)
	return fmt.Sprintf("Response of %s", url), nil
}

func future(url string) <-chan data {
	c := make(chan data, 1)

	go func() {
		body, err := doGet(url)

		c <- data{Body: body, Error: err}
	}()

	return c
}

// используется чтобы запускать в фоне процессы данные которых понадобятся когда-то позднее

func main() {
	future1 := future("https://example1.com")
	future2 := future("https://example2.com")

	fmt.Println("Requests started")

	body1 := <-future1
	body2 := <-future2

	fmt.Println("Response 1:", body1)
	fmt.Println("Response 2:", body2)
}
