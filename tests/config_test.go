package tests

import (
	"github.com/stretchr/testify/assert"
	"go.springy.io/util"
	"testing"
)

func TestConfig(t *testing.T) {
	config := util.Config()
	assert.NotNil(t, config.Database.Name)
	assert.NotNil(t, config.Database.Uri)
	assert.Equal(t, config.Server.Port, 8080)
}
