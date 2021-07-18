package cron

import (
	"fmt"
	"io"
	"reflect"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"

	"github.com/n101661/owl/executors"
)

type config struct {
	skipIfStillRunning  bool
	delayIfStillRunning bool
}

func (cfg config) build() *cron.Cron {
	ws := []cron.JobWrapper{cron.Recover(cron.DefaultLogger)}
	if cfg.skipIfStillRunning {
		ws = append(ws, cron.SkipIfStillRunning(cron.DefaultLogger))
	}
	if cfg.delayIfStillRunning {
		ws = append(ws, cron.DelayIfStillRunning(cron.DefaultLogger))
	}
	return cron.New(cron.WithChain(ws...))
}

// Schedule definition
type Schedule struct {
	name     string
	cron     Cron
	executor executors.Executor
}

// ParseSchedule parses the YAML data as schedules.
func ParseSchedule(r io.Reader) ([]*Schedule, error) {
	dec := yaml.NewDecoder(r)

	var t struct {
		Type string `yaml:"type"`
	}
	if err := dec.Decode(&t); err != nil {
		return nil, err
	}

	builder, ok := getBuilder(t.Type)
	if !ok {
		return nil, fmt.Errorf("unregistered type [%s]", t.Type)
	}

	fields := []reflect.StructField{
		{
			Name: "Name",
			Type: reflect.TypeOf(""),
			Tag:  `yaml:"name"`,
		}, {
			Name: "Cron",
			Type: reflect.TypeOf(Cron{}),
			Tag:  `yaml:"cron"`,
		}, {
			Name: "Config",
			Type: reflect.TypeOf(builder()),
			Tag:  `yaml:"config"`,
		},
	}
	rt := reflect.SliceOf(reflect.StructOf(fields))
	s := reflect.New(rt)
	if err := dec.Decode(s.Interface()); err != nil {
		return nil, err
	}

	s = s.Elem()
	length := s.Len()
	results := make([]*Schedule, length)
	for i := 0; i < length; i++ {
		results[i] = &Schedule{
			name:     s.Index(i).Field(0).String(),
			cron:     s.Index(i).Field(1).Interface().(Cron),
			executor: s.Index(i).Field(2).Interface().(executors.Executor),
		}
	}
	return results, nil
}

func (s *Schedule) cronConfig() config {
	return config{
		skipIfStillRunning:  s.cron.SkipIfStillRunning,
		delayIfStillRunning: s.cron.DelayIfStillRunning,
	}
}

func (s *Schedule) job() cron.Job {
	return cron.FuncJob(func() {
		if err := s.executor.Execute(); err != nil {
			cron.DefaultLogger.Error(err, "finished the job with an error")
		}
		cron.DefaultLogger.Info("finished the job with OK")
	})
}

// Cron definition
type Cron struct {
	Express             string `yaml:"express"`
	SkipIfStillRunning  bool   `yaml:"skip_if_still_running"`
	DelayIfStillRunning bool   `yaml:"delay_if_still_running"`
}

type configBuilder func() executors.Executor
