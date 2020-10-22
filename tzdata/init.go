package tzdata

import (
	"time"
)

func init() {
	// all times should be UTC times!
	utclocation, _ := time.LoadLocation("UTC")
	time.Local = utclocation
}
