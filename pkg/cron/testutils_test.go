package cron

import (
	"fmt"
	"time"

	"github.com/n101661/owl/executors"
)

type testTypeOK struct {
	ID  string `yaml:"id"`
	Age int    `yaml:"age"`
}

func newTestOKBuilder() executors.Executor {
	return &testTypeOK{}
}

func (t *testTypeOK) Execute() error {
	return nil
}

type testTypeFail struct {
	ID  string `yaml:"id"`
	Age int    `yaml:"age"`
}

func newTestFailBuilder() executors.Executor {
	return &testTypeFail{}
}

func (t *testTypeFail) Execute() error {
	return fmt.Errorf("failed")
}

type testTypeLazy struct {
	Sleep int `yaml:"sleep_in_milliseconds"`
}

func newTestLazyBuilder() executors.Executor {
	return &testTypeLazy{}
}

func (t *testTypeLazy) Execute() error {
	time.Sleep(time.Millisecond * time.Duration(t.Sleep))
	return nil
}
