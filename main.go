package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	bufferSize    = 5
	flushInterval = 3 * time.Second
)

type RingBuffer struct {
	data     []int
	size     int
	capacity int
	head     int
	tail     int
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		data:     make([]int, capacity),
		size:     0,
		capacity: capacity,
		head:     0,
		tail:     0,
	}
}

func (rb *RingBuffer) Push(item int) bool {
	if rb.size == rb.capacity {
		return false
	}
	rb.data[rb.tail] = item
	rb.tail = (rb.tail + 1) % rb.capacity
	rb.size++
	return true
}

func (rb *RingBuffer) Pop() (int, bool) {
	if rb.size == 0 {
		return 0, false
	}
	item := rb.data[rb.head]
	rb.head = (rb.head + 1) % rb.capacity
	rb.size--
	return item, true
}

func (rb *RingBuffer) IsEmpty() bool {
	return rb.size == 0
}

func source(out chan<- int) {
	log.Println("Starting source function")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите целые числа (для завершения введите 'q'):")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "q" {
			break
		}
		num, err := strconv.Atoi(input)
		if err == nil {
			out <- num
		} else {
			fmt.Println("Некорректный ввод. Пожалуйста, введите целое число.")
		}
	}
	close(out)
}

func filterNegative(in <-chan int, out chan<- int) {
	log.Println("Starting filterNegative function")
	for num := range in {
		if num >= 0 {
			out <- num
		}
	}
	close(out)
}

func filterNonMultiplesOf3(in <-chan int, out chan<- int) {
	log.Println("Starting filterNonMultiplesOf3 function")
	for num := range in {
		if num != 0 && num%3 == 0 {
			out <- num
		}
	}
	close(out)
}

func bufferStage(in <-chan int, out chan<- int) {
	log.Println("Starting bufferStage function")
	buffer := NewRingBuffer(bufferSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flushBuffer := func() {
		for !buffer.IsEmpty() {
			if item, ok := buffer.Pop(); ok {
				out <- item
			}
		}
	}

	for {
		select {
		case num, ok := <-in:
			if !ok {
				flushBuffer()
				close(out)
				return
			}
			if !buffer.Push(num) {
				flushBuffer()
				buffer.Push(num)
			}
		case <-ticker.C:
			flushBuffer()
		}
	}
}

func consumer(in <-chan int) {
	log.Println("Starting consumer function")
	for num := range in {
		fmt.Printf("Получены данные: %d\n", num)
	}
}

func main() {
	log.Println("Starting main function")
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)

	go source(ch1)
	go filterNegative(ch1, ch2)
	go filterNonMultiplesOf3(ch2, ch3)
	go bufferStage(ch3, ch4)
	consumer(ch4)
}
