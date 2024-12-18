package http

import "gmvr.pw/boombox-web-runner/pkg/model"

type RunnerService interface {
	Run(runner *model.Runner, out chan<- []byte) error
	Stop(runnerId string) (*model.Runner, error)
}
