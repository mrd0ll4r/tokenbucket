package tokenbucket

import (
	"sync"
	"testing"
	"time"
)

func TestConsume(t *testing.T) {
	tb := New(100, 1)
	wg := &sync.WaitGroup{}

	before := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < 10; k++ {
				for !tb.Consume(1) {
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}

	wg.Wait()
	if time.Now().Sub(before) < 1*time.Second {
		t.Fatal("Did not wait 10 seconds")
	}
}

func TestConsumeBurst(t *testing.T) {
	tb := New(1, 10)

	if !tb.Consume(10) {
		t.Fatal("Unable to consume 10 tokens on new bucket with burstSize=10")
	}
}
