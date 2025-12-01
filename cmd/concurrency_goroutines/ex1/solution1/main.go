package main

import "fmt"

func main() {
	workers := 3

	jobs := make(chan int)
	results := make(chan int)
	failed := make(chan int)
	done := make(chan struct{})

	// Start workers
	for i := 0; i < workers; i++ {
		go worker(jobs, results, failed, done)
	}

	// Close results + failed channels after all workers signal done
	go func() {
		for i := 0; i < workers; i++ {
			<-done
		}
		close(results)
		close(failed)
	}()

	// Feed jobs into the queue
	go func() {
		for jobID := 1; jobID <= 20; jobID++ {
			jobs <- jobID
		}
		close(jobs)
	}()

	// Consume results (merged output)
	go func() {
		for r := range results {
			fmt.Printf("Result: %v\n", r)
		}
	}()

	// Main goroutine reads failed jobs
	for f := range failed {
		fmt.Printf("Failed: %v\n", f)
	}

	fmt.Println("All done")
}

func worker(jobs <-chan int, results chan<- int, failed chan<- int, done chan<- struct{}) {
	for job := range jobs {
		success := false

		for attempt := 1; attempt <= 3; attempt++ {
			if job%7 == 0 {
				// fails always for divisible-by-7 (per the exercise)
				continue
			}
			success = true
			results <- job * 10
			break
		}

		if !success {
			failed <- job
		}
	}

	done <- struct{}{}
}
