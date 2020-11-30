package tests

import (
	"github.com/stretchr/testify/assert"
	"go.springy.io/app"
	"testing"
)

// See: http://thylong.com/golang/2016/mocking-mongo-in-golang/
// See: https://medium.com/better-programming/unit-testing-code-using-the-mongo-go-driver-in-golang-7166d1aa72c0
// for mocking database connection
func TestEnv(t *testing.T) {
	env := app.Env()
	assert.NotNil(t, env.Database.Db)
	assert.Equal(t, env.Database.Port, 27017)
	assert.NotNil(t, env.Database.Host)
	assert.NotNil(t, env.Database.Uri())
	assert.Equal(t, env.Server.Port, 8080)
}
