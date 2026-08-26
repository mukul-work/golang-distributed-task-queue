package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mukul-work/golang-distributed-task-queue/internal/handler"
	"github.com/mukul-work/golang-distributed-task-queue/models"
)

func worker(id int, queue <-chan models.Job) {
	for job := range queue {
		fmt.Printf("Worker %d Processing Job %d\n", id, job.Id)
		time.Sleep(5 * time.Second)
	}
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
