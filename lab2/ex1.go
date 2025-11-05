package main

import "fmt"

func main() {
	ch := make(chan string)

	go func() {
		msg := "hey"
		ch <- msg
	}()

	received := <-ch
	fmt.Println("Mesaj primit:", received)
}
