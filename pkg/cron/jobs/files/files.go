package files

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/n101661/owl/pkg/cron"
)

type jobBuilder struct{}

func (b jobBuilder) NewConfig() interface{} {
	return new(Config)
}

type job struct {
	name string
	path string

	values cron.Values
}

func (b jobBuilder) Build(name string, configInterface interface{}) (cron.Job, error) {
	config, ok := configInterface.(*Config)
	if !ok {
		return nil, fmt.Errorf("incompatible config for files job")
	}

	vs := make(cron.Values, len(config.Parameters))
	for _, p := range config.Parameters {
		vs[p.Name] = p.Value
	}
	return &job{
		name:   name,
		path:   config.Path,
		values: vs,
	}, nil
}

func (j *job) Name() string {
	return j.name
}

func (j *job) Values() cron.Values {
	vs := make(cron.Values, len(j.values))
	for k, v := range j.values {
		vs[k] = v
	}
	return vs
}

func (j *job) Run(ctx context.Context, vs cron.Values) error {
	cmd := exec.Command(j.path, parseArgs(vs)...)
	_, err := cmd.Output()
	if e, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf(string(e.Stderr))
	}
	return err
}

func parseArgs(vs cron.Values) []string {
	args, i := make([]string, len(vs)), 0
	for k, v := range vs {
		if v == "" {
			args[i] = k
		} else {
			args[i] = k + "=" + v
		}
		i++
	}
	return args
}

func init() {
	if err := cron.Register("files", jobBuilder{}); err != nil {
		panic("failed to register job builder [files]")
	}
}
