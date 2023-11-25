package mocks

//go:generate mockgen -source ../httputils/interfaces.go -destination httputils/interfaces.go
//go:generate mockgen -source ../app.go -destination app/app.go
//go:generate mockgen -source ../runner.go -destination app/runner.go
