package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/n101661/owl/executors"
	"github.com/n101661/owl/pkg/cron"
)

const builderType = "web hook"

// executor will POST the request with a JSON request.
type executor struct {
	URI        string            `yaml:"uri"`
	Parameters []executors.Param `yaml:"parameters"`
}

func (e *executor) Execute() error {
	req, err := e.parseRequest()
	if err != nil {
		return fmt.Errorf("failed to parse the HTTP request: %v", err)
	}

	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("got a bad status code [%d], response body: %s",
			resp.StatusCode, msg,
		)
	}
	return nil
}

func (e *executor) parseRequest() (*http.Request, error) {
	body := make(map[string]string)
	for _, p := range e.Parameters {
		if _, ok := body[p.Name]; ok {
			continue
		}
		body[p.Name] = p.Value
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, e.URI, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func build() executors.Executor {
	return new(executor)
}

func init() {
	if err := cron.RegisterBuilder(builderType, build); err != nil {
		log.Fatal(err)
	}
}
