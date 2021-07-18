package cron

import "sync"

var runnerRegistry sync.Map

// RegisterBuilder registers a config builder.
func RegisterBuilder(type_ string, builder configBuilder) error {
	if _, loaded := runnerRegistry.LoadOrStore(type_, builder); loaded {
		return ErrTypeExisted
	}
	return nil
}

// ClearRegistry removes all config builders in the
// registry.
func ClearRegistry() {
	runnerRegistry = sync.Map{}
}

func getBuilder(type_ string) (configBuilder, bool) {
	v, ok := runnerRegistry.Load(type_)
	if ok {
		return v.(configBuilder), true
	}
	return nil, false
}
