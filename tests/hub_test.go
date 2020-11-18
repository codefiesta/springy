package tests

import (
	"context"
	"go.springy.io/app"
	"testing"
)

func TestHub(t *testing.T) {

	_, cancel := context.WithTimeout(context.Background(), longDuration)
	defer cancel()

	h1 := app.NewHub()
	h2 := app.NewHub()
	assertEqual(t, h1, h2, "The hub should be a singleton")
}