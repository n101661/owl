package cron

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_registry(t *testing.T) {
	require := require.New(t)

	var (
		builderA        = &testJobBuilder{}
		builderB        = &testJobBuilder{}
		builderAnotherB = &testJobBuilder{}
	)

	require.NoError(Register("a", builderA))
	require.NoError(Register("b", builderB))
	// duplicated `b` type
	require.Error(Register("b", builderAnotherB))

	require.Len(registry, 2)
	require.Same(builderA, registry["a"])
	require.Same(builderB, registry["b"])

	Clear()
	require.Len(registry, 0)
}

func Test_Cron(t *testing.T) {
	require := require.New(t)

	require.NoError(Register("test", &testJobBuilder{
		assertFunc: func(v Values) error {
			require.Len(v, 3)
			require.NotEqual(RandomID, v["a"])

			delete(v, "a")
			require.Equal(Values{
				"b": "b value",
				"c": "c value",
			}, v)
			return nil
		},
	}))

	cron := NewCron()

	// good config
	{
		config := `
---
name: test-job
type: test
cron:
    express: "* * * * *"
    skip_if_still_running: true
    delay_if_still_running: true
---
values:
    a: $RANDOM_ID
    b: b value
    c: c value
`
		require.NoError(cron.AddFromFile(strings.NewReader(config)))
	}
	// bad config
	{
		config := `
---
name: test-job
type: test
cron:
    express: "* * * * *"
    skip_if_still_running: true
    delay_if_still_running: true
`
		require.Error(cron.AddFromFile(strings.NewReader(config)))
	}
	// unknown type
	{
		config := `
---
name: test-job
type: unknown
cron:
    express: "* * * * *"
    skip_if_still_running: true
    delay_if_still_running: true
---
values:
    a: $RANDOM_ID
    b: b value
    c: c value
`
		require.EqualError(
			cron.AddFromFile(strings.NewReader(config)),
			"unknown executor type [unknown]",
		)
	}
}
