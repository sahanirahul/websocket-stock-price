package icore

import "sensibull/stocks-api/business/entities/core"

type IPool interface {
	AddJob(job *core.Job)
}
