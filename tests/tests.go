package tests

import (
	"testing"
	"fmt"
	"time"
)

const shortDuration = 1 * time.Millisecond // a reasonable duration to block in an example

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}