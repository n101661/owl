package cron

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub(t *testing.T) {
	assert := assert.New(t)
	RegisterBuilder("lazy", newTestLazyBuilder)
	defer ClearRegistry()

	const data = `
---
type: lazy
---
- name: example
  cron:
    express: "* * * * *"
    skip_if_still_running: false
    delay_if_still_running: true
  config:
    sleep_in_milliseconds: 100
- name: example
  cron:
    express: "* * * * *"
    skip_if_still_running: true
    delay_if_still_running: false
  config:
    sleep_in_milliseconds: 100`

	// good case
	{
		assert.NotPanics(func() {
			ss, err := ParseSchedule(strings.NewReader(data))
			assert.NoError(err)
			assert.NotEmpty(ss)

			hub := NewHub()
			for i, s := range ss {
				assert.NoErrorf(hub.AddSchedule(s), "at index-%d", i)
			}
			assert.Len(hub.crons, 2)

			assert.NoError(hub.Start())
			assert.NoError(hub.Close(context.Background()))
		})
	}
	// force to close
	{
		assert.NotPanics(func() {
			ss, err := ParseSchedule(strings.NewReader(data))
			assert.NoError(err)
			assert.NotEmpty(ss)

			hub := NewHub()
			for i, s := range ss {
				assert.NoErrorf(hub.AddSchedule(s), "at index-%d", i)
			}

			assert.NoError(hub.Start())

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*20)
			cancel()
			assert.True(errors.Is(hub.Close(ctx), ErrForceClose))
		})
	}
}
