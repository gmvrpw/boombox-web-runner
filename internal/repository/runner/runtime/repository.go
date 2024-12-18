package runtime

import (
	"sync"

	"github.com/google/uuid"
	"gmvr.pw/boombox-web-runner/pkg/model"
)

type RuntimeRunnerRepository struct {
	runners sync.Map
}

func NewRuntimeRunnerRepository() (*RuntimeRunnerRepository, error) {
	return &RuntimeRunnerRepository{runners: sync.Map{}}, nil
}

func (r *RuntimeRunnerRepository) Create(runner *model.Runner) error {
	runner.ID = uuid.New().String()

	r.runners.Store(runner.ID, runner)

	return nil
}

func (r *RuntimeRunnerRepository) DeleteById(id string) (*model.Runner, error) {
	stored, loaded := r.runners.LoadAndDelete(id)
	if !loaded {
		return nil, &model.RunnerNotFoundError{}
	}

	runner, casted := stored.(*model.Runner)
	if !casted {
		return nil, &model.RunnerNotFoundError{}
	}

	return runner, nil
}
