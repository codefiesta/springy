package tests

import (
	"go.springy.io/util"
	"testing"
)

func TestConfig(t *testing.T) {

	c1 := util.Config()
	c2 := util.Config()
	assertEqual(t, c1, c2, "The configuration should be a singleton")
}
