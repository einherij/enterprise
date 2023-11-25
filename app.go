package enterprise

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

type Application interface {
	Run()
	RegisterRunner(runner Runner)
	RegisterOnRun(f func())
	RegisterOnShutdown(f func())
}

type App struct {
	runners []Runner
}

func NewApplication() *App {
	return new(App)
}

func (app *App) RegisterRunner(service Runner) {
	app.runners = append(app.runners, service)
}

func (app *App) RegisterOnRun(f func()) {
	app.runners = append(app.runners, RunnerFunc(f, func() {}))
}

func (app *App) RegisterOnShutdown(f func()) {
	app.runners = append(app.runners, RunnerFunc(func() {}, f))
}

func (app *App) Run() {
	var (
		wg          sync.WaitGroup
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	)
	for _, runnerService := range app.runners {
		serviceRunner := runnerService
		wg.Add(1)
		go func() {
			defer wg.Done()
			serviceRunner.Run(ctx)
		}()
	}

	<-ctx.Done()
	cancel()
	logrus.Info("shutting down application")
	wg.Wait()
}
