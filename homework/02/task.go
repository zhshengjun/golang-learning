package main

import (
	"fmt"
	"sync"
	"time"
)

type Task func()

type TaskResult struct {
	Index    int
	Duration time.Duration
}

func RunTasks(tasks []Task) []TaskResult {
	results := make([]TaskResult, len(tasks))

	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for i, task := range tasks {

		go func(i int, task Task) {
			defer wg.Done()

			start := time.Now()
			task()

			results[i] = TaskResult{
				Index:    i,
				Duration: time.Since(start),
			}
		}(i, task)
	}

	wg.Wait()
	return results
}

func main() {
	tasks := []Task{
		func() { time.Sleep(3 * time.Second) },
		func() { time.Sleep(2 * time.Second) },
		func() { time.Sleep(1 * time.Second) },
	}

	start := time.Now()

	for _, result := range RunTasks(tasks) {
		fmt.Printf("任务 %d：%v\n", result.Index+1, result.Duration)
	}

	fmt.Printf("总耗时：%v\n", time.Since(start))
}
