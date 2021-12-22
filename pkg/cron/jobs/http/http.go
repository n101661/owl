package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cron "github.com/robfig/cron/v3"
)

type jobBuilder struct{}

func (b jobBuilder) NewConfig() interface{} {
	return new(Config)
}

type job struct {
	config Config
}

func (b jobBuilder) Build(configInterface interface{}) (cron.Job, error) {
	config, ok := configInterface.(*Config)
	if !ok {
		return nil, fmt.Errorf("incompatible config for http job")
	}
	return &job{
		config: *config,
	}, nil
}

func (j *job) Run() {
	values := make(map[string]string, len(j.config.Parameters))
	for _, p := range j.config.Parameters {
		values[p.Name] = p.Value
	}

	data, err := json.Marshal(values)
	if err != nil {
		panic(err) // TODO
	}

	resp, err := http.Post(j.config.URI, "application/json", bytes.NewReader(data))
	if err != nil {
		panic(err) // TODO
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		msg := "no response body"
		if resp.ContentLength > 0 {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err) // TODO
			}
			msg = string(data)
		}
		panic(fmt.Errorf("receive bad status code [%d]: %s", code, msg)) // TODO
	}
}
