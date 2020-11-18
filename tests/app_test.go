package tests

import (
	"go.springy.io/app"
	"testing"
)

func TestHub(t *testing.T) {

	h1 := app.NewHub()
	h2 := app.NewHub()
	if h1 != h2 {
		t.Errorf("The hub should be a singleton")
	}
}
