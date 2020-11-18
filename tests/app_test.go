package tests

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {

	//c1 := util.Config()
	//c2 := util.Config()
	//assertEqual(t, c1, c2, "The configuration should be a singleton")
}

func TestHub(t *testing.T) {
	//h1 := app.NewHub()
	//h2 := app.NewHub()
	//assertEqual(t, h1, h2, "The hub should be a singleton")
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}