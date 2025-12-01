package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
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
	Status string // "success" or "failed"
	Value  int
}

var workerCount int32 = 0

func main() {
	topic := "jobs-topic"
	broker := "localhost:9092"
	workerLimit := int32(10)
	maxRetry := 3

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	results := make(chan Result, 100)

	// Kafka writer (produce jobs)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
	defer writer.Close()

	// Seed jobs
	go func() {
		for i := 1; i <= 50; i++ {
			job := Job{ID: i, MaxRetry: maxRetry}
			data, _ := json.Marshal(job)
			err := writer.WriteMessages(ctx, kafka.Message{Value: data})
			if err != nil {
				log.Println("Kafka write error:", err)
			}
		}
		fmt.Println("All jobs produced")
	}()

	// Dynamic worker scaler
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Scale based on current workerCount and limit
				current := atomic.LoadInt32(&workerCount)
				if current < workerLimit {
					atomic.AddInt32(&workerCount, 1)
					go worker(ctx, broker, topic, results)
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Consume results
	go func() {
		for r := range results {
			fmt.Printf("Job %d: %s", r.JobID, r.Status)
			if r.Status == "success" {
				fmt.Printf(", Value=%d", r.Value)
			}
			fmt.Println()
		}
	}()

	// Wait for graceful shutdown
	<-sigs
	fmt.Println("\nShutting down gracefully...")
	cancel()

	// Give workers time to finish
	time.Sleep(3 * time.Second)
	close(results)
	fmt.Println("All done")
}

// worker reads jobs from Kafka and processes them
func worker(ctx context.Context, broker, topic string, results chan<- Result) {
	defer atomic.AddInt32(&workerCount, -1)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  "worker-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.ReadMessage(ctx)
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
				continue
			}

			// Retry logic
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
