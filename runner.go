package enterprise

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type Runner interface {
	Run(ctx context.Context)
}

type runner struct {
	wg  sync.WaitGroup
	log logrus.FieldLogger

	task func(context.Context)
}

func NewRunner(name string, task func(context.Context)) Runner {
	return &runner{
		log:  logrus.WithField("component", name),
		task: task,
	}
}

func (r *runner) run(ctx context.Context) {
	r.log.Info("Run worker")
	defer r.wg.Done()

	r.task(ctx)
}

func (r *runner) Run(ctx context.Context) {
	r.wg.Add(1)
	go r.run(ctx)
}

func RunnerFunc(start, stop func()) Runner {
	return &runnerFunc{start: start, stop: stop}
}

type runnerFunc struct {
	start func()
	stop  func()
}

func (rf *runnerFunc) Run(ctx context.Context) {
	if rf.start != nil {
		rf.start()
	}
	<-ctx.Done()
	if rf.stop != nil {
		rf.stop()
	}
}
