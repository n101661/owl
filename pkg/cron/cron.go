package cron

import (
	"fmt"
	"io"
	"sync"

	cron "github.com/robfig/cron/v3"
	yaml "gopkg.in/yaml.v3"
)

var (
	mu       sync.RWMutex
	registry = map[string]JobBuilder{}
)

func Register(type_ string, builder JobBuilder) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := registry[type_]; ok {
		return fmt.Errorf("duplicated type [%s]", type_)
	}
	registry[type_] = builder
	return nil
}

func Clear() {
	mu.Lock()
	defer mu.Unlock()

	registry = make(map[string]JobBuilder)
}

type JobBuilder interface {
	NewConfig() interface{}
	Build(name string, config interface{}) (Job, error)
}

type Cron struct {
	cron *cron.Cron
}

func NewCron() *Cron {
	return &Cron{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (c *Cron) AddFromFile(r io.Reader) error {
	var cfg Config

	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	mu.RLock()
	builder, ok := registry[cfg.Type]
	mu.RUnlock()
	if !ok {
		return fmt.Errorf("unknown executor type [%s]", cfg.Type)
	}
	jCfg := builder.NewConfig()
	if err := decoder.Decode(jCfg); err != nil {
		return err
	}
	j, err := builder.Build(cfg.Name, jCfg)
	if err != nil {
		return fmt.Errorf("failed to create %s executor: %v", cfg.Name, err)
	}

	job := newJob(j)
	{
		job = withRecover(job)
	}
	if cfg.Cron.SkipIfStillRunning {
		job = withSkipIfStillRunning(job)
	}
	if cfg.Cron.DelayIfStillRunning {
		job = withDelayIfStillRunning(job)
	}

	_, err = c.cron.AddJob(cfg.Cron.Express, job)
	return err
}

func (c *Cron) StartAll() error {
	c.cron.Start()
	return nil
}

func (c *Cron) Close() error {
	<-c.cron.Stop().Done()
	return nil
}
