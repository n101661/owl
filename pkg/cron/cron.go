package cron

import (
	"fmt"
	"io"

	cron "github.com/robfig/cron/v3"
	yaml "gopkg.in/yaml.v3"

	"github.com/n101661/owl/pkg/cron/configs"
)

type JobBuilder interface {
	NewConfig() interface{}
	Build(config interface{}) (cron.Job, error)
}

type Cron struct {
	registry map[string]JobBuilder

	cron *cron.Cron
}

func NewCron() *Cron {
	return &Cron{
		registry: map[string]JobBuilder{},
		cron:     cron.New(),
	}
}

func (c *Cron) Register(type_ string, builder JobBuilder) error {
	if _, ok := c.registry[type_]; ok {
		return fmt.Errorf("duplicated type [%s]", type_)
	}
	c.registry[type_] = builder
	return nil
}

func (c *Cron) AddFromFile(r io.Reader) error {
	var cfg configs.Config

	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	builder, ok := c.registry[cfg.Type]
	if !ok {
		return fmt.Errorf("unknown executor type [%s]", cfg.Type)
	}
	jCfg := builder.NewConfig()
	if err := decoder.Decode(jCfg); err != nil {
		return err
	}
	job, err := builder.Build(jCfg)
	if err != nil {
		return fmt.Errorf("failed to create %s executor: %v", cfg.Name, err)
	}

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
