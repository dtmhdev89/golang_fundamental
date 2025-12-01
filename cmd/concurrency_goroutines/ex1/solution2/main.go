package main

import "fmt"

// Result represents a job's outcome
type Result struct {
	JobID  int
	Status string // "success" or "failed"
	Value  int    // processed value (e.g., job * 10)
}

func main() {
	workers := 3
	jobs := make(chan int)
	results := make(chan Result)
	done := make(chan struct{})

	// Start workers
	for i := 0; i < workers; i++ {
		go worker(jobs, results, done)
	}

	// Close results after all workers finish
	go func() {
		for i := 0; i < workers; i++ {
			<-done
		}
		close(results)
	}()

	// Feed jobs
	go func() {
		for jobID := 1; jobID <= 20; jobID++ {
			jobs <- jobID
		}
		close(jobs)
	}()

	// Read merged results
	for r := range results {
		fmt.Printf("Job %d: %s", r.JobID, r.Status)
		if r.Status == "success" {
			fmt.Printf(", Value=%d", r.Value)
		}
		fmt.Println()
	}

	fmt.Println("All jobs processed")
}

func worker(jobs <-chan int, results chan<- Result, done chan<- struct{}) {
	for job := range jobs {
		success := false
		var value int

		for attempt := 1; attempt <= 3; attempt++ {
			// Simulated failure: divisible by 7 fails
			if job%7 == 0 {
				continue
			}
			success = true
			value = job * 10
			break
		}

		if success {
			results <- Result{JobID: job, Status: "success", Value: value}
		} else {
			results <- Result{JobID: job, Status: "failed"}
		}
	}

	// Signal worker done
	done <- struct{}{}
}
