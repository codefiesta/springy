package util_test

import (
	"github.com/stretchr/testify/assert"
	"go.springy.io/pkg/util"
	"os"
	"path/filepath"
	"testing"
)

// See: http://thylong.com/golang/2016/mocking-mongo-in-golang/
// See: https://medium.com/better-programming/unit-testing-code-using-the-mongo-go-driver-in-golang-7166d1aa72c0
// for mocking database connection
func TestEnv(t *testing.T) {

	// 🤔 The working directory needs to be changed to the root level for this test to run
	wd, _ := os.Getwd()
	root := filepath.Dir(filepath.Dir(wd))
	os.Chdir(root)

	env := util.Env()
	assert.NotNil(t, env.Database.Db)
	assert.Equal(t, env.Database.Port, 27017)
	assert.NotNil(t, env.Database.Host)
	assert.NotNil(t, env.Database.GetURI())
	assert.Equal(t, env.Server.Port, 8080)
}
