//go:build !windows
// +build !windows

package enterprise

import (
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	mock_app "github.com/einherij/enterprise/mocks/app"
)

type ApplicationSuite struct {
	suite.Suite

	app *App

	ctrl       *gomock.Controller
	mockRunner *mock_app.MockRunner
}

func TestWebApplication(t *testing.T) {
	suite.Run(t, new(ApplicationSuite))
}

func (s *ApplicationSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRunner = mock_app.NewMockRunner(s.ctrl)
	s.app = NewApplication()
}

func (s *ApplicationSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ApplicationSuite) TestRegisterRunner() {
	s.Len(s.app.runners, 0)
	s.app.RegisterRunner(s.mockRunner)
	s.Len(s.app.runners, 1)
	s.Equal(s.mockRunner, s.app.runners[0])
}

func (s *ApplicationSuite) TestRegisterOnShutdown() {
	someFunc := func() {}
	s.app.RegisterOnShutdown(someFunc)
	s.Len(s.app.runners, 1)
}

func (s *ApplicationSuite) TestRunAndShutdownBySyscall() {
	serv1 := mock_app.NewMockRunner(s.ctrl)
	serv2 := mock_app.NewMockRunner(s.ctrl)
	s.app.RegisterRunner(serv1)
	s.app.RegisterRunner(serv2)
	serv1.EXPECT().Run(gomock.Any())
	serv2.EXPECT().Run(gomock.Any())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond)
		s.NoError(syscall.Kill(syscall.Getpid(), syscall.SIGINT))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.app.Run()
	}()
	wg.Wait()
}
