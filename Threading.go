package FlowX

import (
	"sync"
)

func SingleThread(fn func()) {
	fn()
}

func MultiThread(fn func(), numThreads int) {
	var wg sync.WaitGroup
	wg.Add(numThreads)

	for i := 0; i < numThreads; i++ {
		fnCopy := fn // Create a copy of fn
		go func() {
			defer wg.Done()
			fnCopy()
		}()
	}

	wg.Wait()
}

func MultiThreadForEachLine(lines []string, fn func(line string), numThreads int) {
	var wg sync.WaitGroup
	wg.Add(len(lines))

	threadPool := make(chan struct{}, numThreads)

	for i := 0; i < numThreads && i < len(lines); i++ {
		line := lines[i]
		threadPool <- struct{}{}
		fnCopy := fn // Create a copy of fn
		go func(line string) {
			defer func() { <-threadPool }()
			defer wg.Done()
			fnCopy(line)
		}(line)
	}

	wg.Wait()
}
