package faults

import (
	"fmt"
	"time"
)

func Retry(operation func() error, maxRetries, baseDelay int) error {
	delay := time.Millisecond * time.Duration(baseDelay)
	for n := 0; n < maxRetries; n++ {
		if nil == operation() {
			return nil
		}
		time.Sleep(delay)
		delay = delay * 2
	}
	return fmt.Errorf("retry failed after %d attempt(s) with total delay: %v s", maxRetries, delay)
}

func TimeOut(operation func() error, timeout int) error {
	timer := time.NewTimer(time.Duration(timeout) * time.Millisecond)
	done := make(chan error)
	defer close(done)
	go func() {
		done <- operation()
	}()
	select {
	case <-timer.C:
		return fmt.Errorf("timeout exceeded")
	case err := <-done:
		if err != nil {
			return err
		} else {
			return nil
		}

	}
}

type DeadLetter struct {
	message string
	err     error
}

type DeadLetterQueue struct {
	Queue []DeadLetter
	len   int
}

func NewDeadLetterQueue(capacity int) DeadLetterQueue {
	var dlq DeadLetterQueue
	dlq.Queue = make([]DeadLetter, 1, capacity)
	dlq.len = 1
	return dlq
}

func (dlq *DeadLetterQueue) append(msg string, err error) {
	if dlq.len == cap(dlq.Queue) {
		newQueue := make([]DeadLetter, dlq.len, 2*cap(dlq.Queue))
		copy(newQueue, dlq.Queue)
		dlq.Queue = newQueue
	}
	newDL := &(dlq.Queue[dlq.len-1])
	newDL.message = msg
	newDL.err = err
	dlq.len += 1
}

func (dlq *DeadLetterQueue) GetMessages() []string {
	messages := make([]string, dlq.len)
	for i, dl := range dlq.Queue {
		messages[i] = dl.message
		fmt.Println(dl.message)
	}
	return messages
}

func ProcessWithDLQ(messages []string, handler func(string) error, dlq DeadLetterQueue) {
	for _, msg := range messages {
		err := handler(msg)
		if err != nil {
			dlq.append(msg, err)
		}
	}
}
