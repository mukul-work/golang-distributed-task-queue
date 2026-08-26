package models

type Job struct {
	Id          int    `json:"id"`
	Task        string `json:"task"`
	Attempts    int
	MaxAttempts int `json:"maxAttempts"`
}
