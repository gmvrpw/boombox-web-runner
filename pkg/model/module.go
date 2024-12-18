package model

type Selector = string

type Action struct {
	Selector Selector `yaml:"selector"`
	Humanly  bool     `yaml:"humanly"`
}

type Time struct {
	Selector Selector `yaml:"selector"`
}

type Module struct {
	Name string `yaml:"name"`
	Test string `yaml:"test"`
	Auth struct {
		Name   string `yaml:"name"`
		Value  string `yaml:"value"`
		Domain string `yaml:"domain"`
	} `yaml:"auth"`
	PlaybackTime  Time   `yaml:"playback"`
	DurationTime  Time   `yaml:"duration"`
	RemainingTime Time   `yaml:"remaining"`
	PlayAction    Action `yaml:"play"`
}

type ModuleNotFoundError struct{}

func (e *ModuleNotFoundError) Error() string {
	return "module not found"
}
