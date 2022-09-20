package asyncjob

import (
	"context"
	"food-delivery/common"
	"log"
	"sync"
)

type group struct {
	isConcurrent bool
	jobs         []Job
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	g := &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}

	return g
}

func (g *group) Run(context context.Context) error {
	g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))

	for i, _ := range g.jobs {
		if g.isConcurrent {
			go func(aj Job) {
				defer common.AppRecover()
				errChan <- g.runJob(context, aj)
				g.wg.Done()
			}(g.jobs[i])
			continue
		}
		job := g.jobs[i]
		errChan <- g.runJob(context, job)
		g.wg.Done()
	}
	var err error
	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
		}
	}
	g.wg.Wait()
	return err
}

func (g *group) runJob(context context.Context, j Job) error {
	if err := j.Execute(context); err != nil {
		for {
			log.Println(err)
			if j.State() == StateRetryFailed {
				return err
			}
			if j.Retry(context) == nil {
				return nil
			}

		}
	}
	return nil
}
