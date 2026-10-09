package valkey

import (
	"time"
)

const (
	CacheQueryTimeout = 500 * time.Millisecond

	UserExpTime    = time.Hour
	TeacherExpTime = time.Hour
	StudentExpTime = time.Hour
)
