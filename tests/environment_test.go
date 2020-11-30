package tests

import (
	"github.com/stretchr/testify/assert"
	"go.springy.io/app"
	"testing"
)

func TestEnv(t *testing.T) {
	env := app.Env()
	assert.NotNil(t, env.Database.Db)
	assert.Equal(t, env.Database.Port, 27017)
	assert.NotNil(t, env.Database.Host)
	assert.NotNil(t, env.Database.Uri())
	assert.Equal(t, env.Server.Port, 8080)
}
