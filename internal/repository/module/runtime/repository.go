package runtime

import (
	"regexp"

	"gmvr.pw/boombox-web-runner/pkg/model"
)

type RuntimeModuleRepostory struct {
	modules []model.Module
}

func NewRuntimeModuleRepository(modules []model.Module) (*RuntimeModuleRepostory, error) {
	return &RuntimeModuleRepostory{modules: modules}, nil
}

func (r *RuntimeModuleRepostory) GetModuleByUrl(url string) (*model.Module, error) {
	for _, module := range r.modules {
		matched, err := regexp.MatchString(module.Test, url)
		if err != nil {
			return nil, err
		}

		if matched {
			m := module
			return &m, nil
		}
	}
	return nil, &model.ModuleNotFoundError{}
}
