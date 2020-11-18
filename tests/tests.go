package tests

import (
	"testing"
	"fmt"
	"time"
)

const shortDuration = 1 * time.Millisecond // a reasonable duration to block in a test
const longDuration = 1 * time.Second // the longest duration to block in a test

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func assertNotNil(t *testing.T, a interface{}, message string) {
	if a != nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != nil", a)
	}
	t.Fatal(message)
}