package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/n101661/owl/pkg/cron"
	"github.com/stretchr/testify/require"
)

func Test_jobBuilder_NewConfig(t *testing.T) {
	require := require.New(t)

	builder := jobBuilder{}

	require.Equal(new(Config), builder.NewConfig())
}

func Test_jobBuilder_Build(t *testing.T) {
	require := require.New(t)

	builder := jobBuilder{}

	// good case
	{
		serverOK := true
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if serverOK {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("gg"))
		}))
		defer s.Close()

		j, err := builder.Build("test", &Config{
			URI: s.URL,
			Parameters: []cron.Parameter{
				{Name: "id", Value: cron.RandomID},
				{Name: "a", Value: "aaa"},
				{Name: "b", Value: "bbb"},
			},
		})
		require.NoError(err)
		require.Equal(&job{
			name: "test",
			uri:  s.URL,
			values: cron.Values{
				"id": cron.RandomID,
				"a":  "aaa",
				"b":  "bbb",
			},
		}, j)

		require.Equal("test", j.Name())
		require.Equal(cron.Values{
			"id": cron.RandomID,
			"a":  "aaa",
			"b":  "bbb",
		}, j.Values())
		// call server ok
		{
			require.NoError(j.Run(context.Background(), j.Values()))
		}
		// call server failed
		{
			serverOK = false
			require.EqualError(j.Run(context.Background(), j.Values()), "receive bad status code [400]: gg")
		}
	}
	// incompatible config
	{
		j, err := builder.Build("test", struct{}{})
		require.Error(err)
		require.Nil(j)
	}
}
