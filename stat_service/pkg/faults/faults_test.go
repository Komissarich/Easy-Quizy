package faults_test

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"quiz_app/pkg/faults"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	operation := func() error {
		ok := rand.Float32()
		if ok > 0.999 {
			return nil
		} else {
			return fmt.Errorf("internal operation error")
		}
	}
	if faults.Retry(operation, 5, 100) == nil {
		t.Errorf("retry test failed")
	}
}

func TestTimeOut(t *testing.T) {
	operation := func() error {
		ok := rand.Float32()
		if ok < 0.7 {
			time.Sleep(2 * time.Second)
			return nil
		} else {
			return fmt.Errorf("internal operation error")
		}
	}

	if faults.TimeOut(operation, 1000) == nil {
		t.Errorf("timeout test failed")
	}
}

func TestProccessWithDLQ(t *testing.T) {
	messages := []string{"msg1", "msg2", "msg3"}
	dlq := faults.NewDeadLetterQueue(3)

	faults.ProcessWithDLQ(messages, func(msg string) error {
		if msg == "msg2" {
			return errors.New("processing failed")
		}
		// fmt.Printf("Processed: %s\n", msg)
		return nil
	}, dlq)

	// fmt.Println(dlq.GetMessages())
	if len(dlq.GetMessages()) != 1 {
		t.Errorf("dead letter queue test failed")
	}

}
