package http

type Config struct {
	URI        string      `yaml:"uri"`
	Parameters []Parameter `yaml:"parameters"`
}

type Parameter struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
