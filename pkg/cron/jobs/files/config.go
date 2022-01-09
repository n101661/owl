package files

import "github.com/n101661/owl/pkg/cron"

type Config struct {
	Path       string           `yaml:"path"`
	Parameters []cron.Parameter `yaml:"parameters"`
}
