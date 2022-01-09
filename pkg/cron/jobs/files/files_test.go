package files

import (
	"testing"

	"github.com/n101661/owl/pkg/cron"
	"github.com/stretchr/testify/require"
)

func Test_parseArgs(t *testing.T) {
	require := require.New(t)

	// good case
	{
		args := parseArgs(cron.Values{
			"--name":             "tester",
			"--value-with-space": "hello world",
			"--ok":               "",
		})
		require.ElementsMatch([]string{
			"--name=tester",
			"--value-with-space=hello world",
			"--ok",
		}, args)
	}
}
