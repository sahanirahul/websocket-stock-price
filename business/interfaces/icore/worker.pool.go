package icore

import "priceupdater/stocks-api/business/entities/core"

type IPool interface {
	AddJob(job *core.Job)
}
