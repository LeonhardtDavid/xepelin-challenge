package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func Test_LoadConfig_Success(t *testing.T) {
	expectedPort := 9999
	os.Setenv("APP_PORT", strconv.Itoa(expectedPort))

	config, err := LoadConfig()

	os.Clearenv()

	assert.Nil(t, err)
	assert.Equal(t, expectedPort, config.Port)
}

func Test_LoadConfig_Defaults(t *testing.T) {
	config, err := LoadConfig()

	assert.Nil(t, err)
	assert.Equal(t, 8080, config.Port)
}

func Test_LoadConfig_PortEnvNotInt(t *testing.T) {
	os.Setenv("APP_PORT", "not_an_int")

	config, err := LoadConfig()

	os.Clearenv()

	assert.Nil(t, err)
	assert.Equal(t, 8080, config.Port)
}

func Test_LoadConfig_InvalidPortEnv(t *testing.T) {
	os.Setenv("APP_PORT", "-1")

	config, err := LoadConfig()

	os.Clearenv()

	assert.Error(t, err)
	assert.Nil(t, config)
}
