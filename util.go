package main 

import "time"

var now time.Time

func getNow() time.Time {
	if now.IsZero() {
		return time.Now()
	} else {
		return now
	}
}
