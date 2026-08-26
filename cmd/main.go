package main

import (
	"fmt"
	"net/http"

	"github.com/mukul-work/golang-distributed-task-queue/internal/handler"
	"github.com/mukul-work/golang-distributed-task-queue/models"
)

func worker(queue <-chan models.Job) {
	for job := range queue {
		fmt.Printf("Processing Job %d", job.Id)
	}
}

func main() {
	queue := make(chan models.Job, 10)
	h := handler.Handler{
		Queue: queue,
	}

	go worker(queue)

	mux := http.NewServeMux()
	mux.HandleFunc("/create", h.CreateJob)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)

}
