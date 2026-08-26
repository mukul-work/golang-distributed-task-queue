package main

import (
	"fmt"
	"time"

	"github.com/mukul-work/golang-distributed-task-queue/models"
)

func worker(queue <-chan models.Job) {
	for job := range queue {
		fmt.Printf("Processing Job %d", job.Id)
	}
}

func main() {
	queue := make(chan models.Job, 10)

	go worker(queue)

	job := models.Job{
		Id:   1,
		Task: "Do nothing",
	}

	queue <- job

	time.Sleep(time.Second)

}
