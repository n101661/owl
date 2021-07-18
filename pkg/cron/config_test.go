package cron

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSchedule(t *testing.T) {
	assert := assert.New(t)
	const typeTest = "test"

	assert.NoError(RegisterBuilder(typeTest, newTestOKBuilder))
	defer ClearRegistry()

	// good case
	{
		const data = `
---
type: test
---
- name: example
  cron:
    express: "5 * * * *"
    skip_if_still_running: false
    delay_if_still_running: true
  config:
    id: tester
    age: 17`
		s, err := ParseSchedule(strings.NewReader(data))
		assert.NoError(err)
		assert.Len(s, 1)
		assert.Equal("example", s[0].name)
		assert.Equal(Cron{
			Express:             "5 * * * *",
			SkipIfStillRunning:  false,
			DelayIfStillRunning: true,
		}, s[0].cron)
		assert.NoError(s[0].executor.Execute())
	}
	// unknown type
	{
		const data = `
---
type: unknown
---
- name: example
  cron:
    express: "5 * * * *"
    skip_if_still_running: false
    delay_if_still_running: true
  config:
    id: tester
    age: 17`
		s, err := ParseSchedule(strings.NewReader(data))
		assert.EqualError(err, "unregistered type [unknown]")
		assert.Empty(s)
	}
}
