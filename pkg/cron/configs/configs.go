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

type HTTP struct {
	URI        string      `yaml:"uri"`
	Parameters []HTTPParam `yaml:"parameters"`
}

type HTTPParam struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
