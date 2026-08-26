package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mukul-work/golang-distributed-task-queue/internal/handler"
	"github.com/mukul-work/golang-distributed-task-queue/models"
)

func worker(id int, queue <-chan models.Job) {
	for job := range queue {
		fmt.Printf("Worker %d picked Job %d\n", id, job.Id)
		err := process(job)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Worker %d completed Job %d\n\n", id, job.Id)
	}
}

func process(job models.Job) error {
	job.Attempts++
	fmt.Printf("Attempt no. : %d\n", job.Attempts)
	if job.MaxAttempts < job.Attempts {
		return errors.New("Job removed from queue\n")
	}
	if job.Id == 2 {
		return errors.New("Job failed\n")
	}

	return nil
}

func main() {
	queue := make(chan models.Job, 10)
	h := handler.Handler{
		Queue: queue,
	}

	// start 4 workers
	for i := 0; i < 4; i++ {
		go worker(i, queue)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/create", h.CreateJob)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)

}
