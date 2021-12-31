package cron

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Values map[string]string `yaml:"values"`

	assertFunc func(Values) error `yaml:"-"`
}

type testJobBuilder struct {
	assertFunc func(Values) error
}

func (b *testJobBuilder) NewConfig() interface{} {
	return &testConfig{
		assertFunc: b.assertFunc,
	}
}

func (b *testJobBuilder) Build(name string, config interface{}) (Job, error) {
	cfg := config.(*testConfig)
	return &testJob{
		values:     cfg.Values,
		assertFunc: cfg.assertFunc,
	}, nil
}

type testJob struct {
	values Values
	// assertFunc it will assert values when Run is called.
	assertFunc func(Values) error
}

func (j *testJob) Name() string { return "test job" }

func (j *testJob) Values() Values {
	return j.values
}

func (j *testJob) Run(ctx context.Context, vs Values) error {
	return j.assertFunc(vs)
}

func Test_newJob(t *testing.T) {
	require := require.New(t)

	const (
		keyHa   = "hahaha"
		keyName = "name"
	)

	j := &testJob{
		values: Values{
			keyHa:   RandomID,
			keyName: "tester",
		},
		assertFunc: func(v Values) error {
			require.Len(v, 2)
			require.NotEqual(RandomID, v[keyHa])

			delete(v, keyHa)
			require.Equal(Values{keyName: "tester"}, v)
			return nil
		},
	}

	job := newJob(j)
	require.NotPanics(func() {
		job.Run()
	})
}

func Test_newValuesCache(t *testing.T) {
	require := require.New(t)

	// values with reserved value
	{
		cache := newValuesCache(Values{
			"no":     "123",
			"hahaha": RandomID,
		})
		require.Equal(&valuesCache{
			randomIDKey: []string{"hahaha"},
		}, cache)
	}
	// values without reserved value
	{
		cache := newValuesCache(Values{
			"no":     "123",
			"hahaha": "RANDOM_ID",
		})
		require.Equal(&valuesCache{
			randomIDKey: []string{},
		}, cache)
	}
}
