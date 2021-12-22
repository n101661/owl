package configs

type Config struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Cron Cron   `yaml:"cron"`
}

type Cron struct {
	Express             string `yaml:"express"`
	SkipIfStillRunning  bool   `yaml:"skip_if_still_running"`
	DelayIfStillRunning bool   `yaml:"delay_if_still_running"`
}
