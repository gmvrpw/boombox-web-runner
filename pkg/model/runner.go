package model

type Runner struct {
	ID       string `json:"id"`
	Url      string `json:"url"`
	Timecode uint64 `json:"timecode,omitempty"`
	Port     int    `json:"port"`

	// only if recieved from repository
	Stop chan<- bool `json:"-"`
}

type RunnerNotFoundError struct{}

func (e *RunnerNotFoundError) Error() string {
	return "runner not found"
}
