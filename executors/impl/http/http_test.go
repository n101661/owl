package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/n101661/owl/executors"
)

type testServer struct {
	values map[string]string
}

func newTestServer(vs map[string]string) *testServer {
	return &testServer{values: vs}
}

func (s *testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	actual, err := s.getValues(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	if !reflect.DeepEqual(actual, s.values) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("expected values: %v, got: %v", s.values, actual)))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *testServer) getValues(r *http.Request) (map[string]string, error) {
	vs := make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&vs); err != nil {
		return nil, err
	}
	return vs, nil
}

func TestExecutor(t *testing.T) {
	assert := assert.New(t)

	// good case
	{
		s := httptest.NewServer(newTestServer(map[string]string{
			"id":   "123",
			"name": "a",
		}))

		e := executor{
			URI: s.URL,
			Parameters: []executors.Param{
				{Name: "id", Value: "123"},
				{Name: "name", Value: "a"},
			},
		}
		assert.NoError(e.Execute())
	}
	// missing parameters
	{
		s := httptest.NewServer(newTestServer(map[string]string{}))

		e := executor{
			URI:        s.URL,
			Parameters: []executors.Param{},
		}
		assert.NoError(e.Execute())
	}
	// duplicated parameter
	{
		s := httptest.NewServer(newTestServer(map[string]string{
			"id":   "123",
			"name": "get first",
		}))

		e := executor{
			URI: s.URL,
			Parameters: []executors.Param{
				{Name: "id", Value: "123"},
				{Name: "name", Value: "get first"},
				{Name: "name", Value: "ignore second"},
				{Name: "name", Value: "ignore others"},
			},
		}
		assert.NoError(e.Execute())
	}
	// bad request
	{
		s := httptest.NewServer(newTestServer(map[string]string{
			"id":   "123",
			"name": "a",
		}))

		e := executor{
			URI:        s.URL,
			Parameters: []executors.Param{},
		}
		assert.Error(e.Execute())
	}
}
