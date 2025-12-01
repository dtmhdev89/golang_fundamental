package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

type Job struct {
	ID       int
	Attempts int
	MaxRetry int
}

type Result struct {
	JobID  int
	Status string
	Value  int
}

func main() {
	topic := "jobs-topic"
	broker := "localhost:9092"
	workerLimit := 10
	maxRetry := 3

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Produce jobs
	go produceJobs(ctx, broker, topic, maxRetry)

	// Create worker pool with job channel
	jobs := make(chan Job, 100)
	results := make(chan Result, 100)

	// Start fixed number of workers
	var wg sync.WaitGroup
	for i := 0; i < workerLimit; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i+1, jobs, results)
	}

	// Single Kafka consumer feeding the worker pool
	go consumeJobs(ctx, broker, topic, jobs)

	// Print results
	go func() {
		for r := range results {
			fmt.Printf("Job %d: %s", r.JobID, r.Status)
			if r.Status == "success" {
				fmt.Printf(", Value=%d", r.Value)
			}
			fmt.Println()
		}
	}()

	// Wait for shutdown signal
	<-sigs
	fmt.Println("\nShutting down gracefully...")
	cancel()

	// Close jobs channel and wait for workers to finish
	close(jobs)
	wg.Wait()
	close(results)

	fmt.Println("All done")
}

func produceJobs(ctx context.Context, broker, topic string, maxRetry int) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker},
		Topic:        topic,
		BatchSize:    10, // Batch messages
		BatchTimeout: 10 * time.Millisecond,
		Async:        true, // Don't wait for acknowledgment
	})
	defer writer.Close()

	// Create batch of messages
	messages := make([]kafka.Message, 50)
	for i := 0; i < 50; i++ {
		job := Job{ID: i + 1, MaxRetry: maxRetry}
		data, _ := json.Marshal(job)
		messages[i] = kafka.Message{Value: data}
	}

	// Write all messages in batch
	err := writer.WriteMessages(ctx, messages...)
	if err != nil {
		log.Println("Kafka write error:", err)
		return
	}
	fmt.Println("All 50 jobs produced")
}

func consumeJobs(ctx context.Context, broker, topic string, jobs chan<- Job) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  "worker-group",
		MinBytes: 1, // Don't wait for data to accumulate
		MaxBytes: 10e6,
		MaxWait:  100 * time.Millisecond, // Max wait time for batch
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Println("Kafka read error:", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			var job Job
			if err := json.Unmarshal(msg.Value, &job); err != nil {
				log.Println("Invalid job format:", err)
				reader.CommitMessages(ctx, msg)
				continue
			}

			// Send job to worker pool
			select {
			case jobs <- job:
				// Commit after successfully queuing
				reader.CommitMessages(ctx, msg)
			case <-ctx.Done():
				return
			}
		}
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, id int, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			// Process job with retry logic
			success := false
			for job.Attempts < job.MaxRetry {
				job.Attempts++
				if job.ID%7 != 0 { // simulate failure for divisible by 7
					success = true
					break
				}
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
