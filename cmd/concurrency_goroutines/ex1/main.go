// Goal
// Implement a worker pool where:
// You have N workers.
// They receive jobs (integers representing "task IDs").
// Each job must be processed with retry up to 3 times.
// If a job fails 3 times, send it to a failedJobs channel.
// Otherwise, send the job result (jobID × 10) to results.
// Constraints
// Do not use sync.WaitGroup.
// You must coordinate worker termination via channel closing only.
// Workers must exit gracefully when job channel closes.
// Main goroutine must wait until:
// All workers have exited
// All results & failedJobs have been processed
// Required Output
// Print processed results in order they finish.
// Print failed jobs at the end.
// Inputs
// Jobs: 1..20
// Retry rule:
// Fail processing if jobID is divisible by 7 (7, 14).
// Hints
// Use at least three channels:
// jobs, results, failedJobs
// Implement retries inside each worker.
// Gracefully close results and failedJobs after workers complete.

package main

import "fmt"

func main() {
	workers := int(3)

	jobs := make(chan int)
	results := make(chan int)
	failedJobs := make(chan int)
	done := make(chan struct{})

	for i := 0; i < workers; i++ {
		go processJobs(jobs, results, failedJobs, done)
	}

	go func() {
		for i := 0; i < workers; i++ {
			<-done
		}

		close(results)
		close(failedJobs)
	}()

	go func() {
		for jobID := range 20 {
			jobs <- jobID + 1
		}

		close(jobs)
	}()

	go func() {
		for v := range results {
			fmt.Printf("Results: %v\n", v)
		}
	}()

	for v := range failedJobs {
		fmt.Printf("failedJobs: %v\n", v)
	}
}

func processJobs(jobs chan int, results chan int, failedJob chan int, done chan struct{}) {
	for jobID := range jobs {
		failed := false
		for range 3 {
			failed = false
			if jobID%7 == 0 {
				failed = true
				continue
			}

			results <- jobID
			break
		}
		if failed {
			failedJob <- jobID
		}
	}

	done <- struct{}{}
}
