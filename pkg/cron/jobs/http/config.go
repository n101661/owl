package http

import "github.com/n101661/owl/pkg/cron"

type Config struct {
	URI        string           `yaml:"uri"`
	Parameters []cron.Parameter `yaml:"parameters"`
}
