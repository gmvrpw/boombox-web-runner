package runner

import "gmvr.pw/boombox-web-runner/pkg/model"

type RunnerRepository interface {
	Create(*model.Runner) error

	DeleteById(id string) (*model.Runner, error)
}

type ModuleRepository interface {
	GetModuleByUrl(url string) (*model.Module, error)
}
