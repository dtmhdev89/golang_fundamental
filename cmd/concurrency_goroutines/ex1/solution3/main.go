package main

import (
	"context"
	"fmt"
	"time"
)

type Result struct {
	JobID  int
	Status string // "success" or "failed"
	Value  int
}

type Job struct {
	ID       int
	Attempts int
	MaxRetry int
}

func main() {
	workerCount := 5
	jobCount := 50
	maxRetry := 3

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan Job, 100) // buffered for scalability
	results := make(chan Result, 100)

	// Start workers
	for i := 0; i < workerCount; i++ {
		go worker(ctx, jobs, results)
	}

	// Produce jobs
	go func() {
		for i := 1; i <= jobCount; i++ {
			jobs <- Job{ID: i, Attempts: 0, MaxRetry: maxRetry}
		}
		close(jobs)
	}()

	// Consume results
	for i := 0; i < jobCount; i++ {
		res := <-results
		fmt.Printf("Job %d: %s", res.JobID, res.Status)
		if res.Status == "success" {
			fmt.Printf(", Value=%d", res.Value)
		}
		fmt.Println()
	}

	fmt.Println("All jobs processed")
}

func worker(ctx context.Context, jobs <-chan Job, results chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			// Retry loop
			var success bool
			for job.Attempts < job.MaxRetry {
				job.Attempts++
				if processJob(job.ID) {
					success = true
					break
				}
				// Optional backoff
				time.Sleep(time.Duration(job.Attempts) * 100 * time.Millisecond)
			}

			if success {
				results <- Result{JobID: job.ID, Status: "success", Value: job.ID * 10}
			} else {
				results <- Result{JobID: job.ID, Status: "failed"}
			}
		}
	}
}

// Simulate job processing: divisible by 7 fails
func processJob(jobID int) bool {
	if jobID%7 == 0 {
		return false
	}
	return true
}
