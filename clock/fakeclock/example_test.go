package fakeclock

import (
	"fmt"
	"sync"
	"time"
)

func ExampleClock_Next() {
	var c Clock

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		i := i
		timer := c.Timer(time.Duration(i+1) * time.Second)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer timer.Stop()

			now := <-timer.C()
			fmt.Printf("timer %d fired at %v\n", i, now)
		}()
	}
	for {
		if _, ok := c.Next(); !ok {
			break
		}
	}
	wg.Wait()

	// Unordered output:
	// timer 0 fired at 0001-01-01 00:00:01 +0000 UTC
	// timer 1 fired at 0001-01-01 00:00:02 +0000 UTC
	// timer 2 fired at 0001-01-01 00:00:03 +0000 UTC
}
