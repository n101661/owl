package cron

type Config struct {
	Name string     `yaml:"name"`
	Type string     `yaml:"type"`
	Cron CronConfig `yaml:"cron"`
}

type CronConfig struct {
	Express             string `yaml:"express"`
	SkipIfStillRunning  bool   `yaml:"skip_if_still_running"`
	DelayIfStillRunning bool   `yaml:"delay_if_still_running"`
}

type Parameter struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
