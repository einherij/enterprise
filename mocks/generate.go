package mocks

//go:generate mockgen -source ../app.go -destination app/app.go
//go:generate mockgen -source ../runner/runner.go -destination runner/runner.go
//go:generate mockgen -source ../httputils/interfaces.go -destination httputils/interfaces.go
