package queue

import "unify-backend/internal/job"

var ADBQueue chan job.ADBJob

func InitADBQueue(buffer int) {
	ADBQueue = make(chan job.ADBJob, buffer)
}

func EnqueueADB(j job.ADBJob) {
	ADBQueue <- j
}

func GetADBQueue() chan job.ADBJob {
	return ADBQueue
}