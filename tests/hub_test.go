package tests

import (
	"go.springy.io/app"
	"testing"
)

func TestHub(t *testing.T) {

	hub := app.NewHub()
	assertNotNil(t, hub, "The hub should not be nil")
}