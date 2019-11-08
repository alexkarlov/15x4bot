package time

import (
	"time"
)

var t func() time.Time

func init() {
	t = time.Now
}

func SetNow(n func() time.Time) {
	t = n
}

func Now() time.Time {
	return t()
}
