package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/n101661/owl/pkg/cron"
)

type jobBuilder struct{}

func (b jobBuilder) NewConfig() interface{} {
	return new(Config)
}

type job struct {
	name string
	uri  string

	values cron.Values
}

func (b jobBuilder) Build(name string, configInterface interface{}) (cron.Job, error) {
	config, ok := configInterface.(*Config)
	if !ok {
		return nil, fmt.Errorf("incompatible config for http job")
	}

	vs := make(cron.Values, len(config.Parameters))
	for _, p := range config.Parameters {
		vs[p.Name] = p.Value
	}
	return &job{
		name:   name,
		uri:    config.URI,
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
	data, err := json.Marshal(vs)
	if err != nil {
		return err
	}

	resp, err := http.Post(j.uri, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		msg := "no response body"
		if resp.ContentLength > 0 {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			msg = string(data)
		}
		return fmt.Errorf("receive bad status code [%d]: %s", code, msg)
	}
	return nil
}

func init() {
	if err := cron.Register("http", jobBuilder{}); err != nil {
		panic("failed to register job builder [http]")
	}
}
