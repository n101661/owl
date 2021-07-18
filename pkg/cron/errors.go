package cron

import "fmt"

var (
	ErrForceClose  = fmt.Errorf("force to close the hub")
	ErrTypeExisted = fmt.Errorf("the type is already existed")
)
