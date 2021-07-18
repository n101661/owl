package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterBuilder(t *testing.T) {
	assert := assert.New(t)

	const (
		typeOK   = "ok"
		typeFail = "fail"
	)
	assertRegistry := func(name string, expectedTypes ...string) {
		types := []string{}
		runnerRegistry.Range(func(key, _ interface{}) bool {
			s, _ := key.(string)
			types = append(types, s)
			return true
		})
		assert.Lenf(types, len(expectedTypes), "case [%s]", name)
		for _, t := range expectedTypes {
			assert.Containsf(types, t, "case [%s]: missing %s", name, t)
		}
	}

	// first register
	assert.NoError(RegisterBuilder(typeOK, newTestOKBuilder))
	assertRegistry("first register", typeOK)

	// second register
	assert.NoError(RegisterBuilder(typeFail, newTestFailBuilder))
	assertRegistry("second register", typeOK, typeFail)

	// register with duplicated type
	assert.Error(RegisterBuilder(typeOK, newTestOKBuilder))
	assertRegistry("register with duplicated type", typeOK, typeFail)

	ClearRegistry()
	assertRegistry("clear")
}

func Test_getBuilder(t *testing.T) {
	const (
		typeOK      = "ok"
		typeUnknown = "??"
	)
	assert := assert.New(t)
	assert.NoError(RegisterBuilder(typeOK, newTestOKBuilder))

	// get a builder
	{
		builder, ok := getBuilder(typeOK)
		assert.True(ok)
		assert.NotNil(builder)
	}
	// builder not found
	{
		builder, ok := getBuilder(typeUnknown)
		assert.False(ok)
		assert.Nil(builder)
	}
}
