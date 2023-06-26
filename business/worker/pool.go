package worker

import (
	"fmt"
	"sensibull/stocks-api/business/entities/core"
	"sensibull/stocks-api/business/interfaces/icore"
)

type worker struct {
	numberOfWorker int
	jobs           chan *core.Job
}

func (w *worker) AddJob(job *core.Job) {
	w.jobs <- job
}

func (w *worker) start() {
	fmt.Println("starting workers")
	for i := 1; i <= w.numberOfWorker; i++ {
		go w.run(i)
	}
}

func (w *worker) run(workerId int) {
	fmt.Println("starting worker ", workerId)
	for {
		job := <-w.jobs
		job.Execute()
	}
}

func NewWorkerPool(maxworkers int, jobQSize int) icore.IPool {
	w := new(worker)
	w.jobs = make(chan *core.Job, jobQSize)
	w.numberOfWorker = maxworkers
	w.start()
	return w
}
