package tests

import (
	"github.com/ip-05/quizzus/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigServer(t *testing.T) {
	cfg := config.InitConfig("test")

	assert.Equal(t, false, cfg.Server.Secure, "should be equal")
	assert.Equal(t, "localhost", cfg.Server.Domain, "should be equal")
	assert.Equal(t, "http://localhost:1234", cfg.Server.Base, "should be equal")
	assert.Equal(t, "localhost", cfg.Server.Host, "should be equal")
	assert.Equal(t, int64(1234), cfg.Server.Port, "should be equal")
}

func TestConfigGoogle(t *testing.T) {
	cfg := config.InitConfig("test")

	assert.Equal(t, "id", cfg.Google.ClientId, "should be equal")
	assert.Equal(t, "secret", cfg.Google.ClientSecret, "should be equal")
}

func TestConfigSecrets(t *testing.T) {
	cfg := config.InitConfig("test")

	assert.Equal(t, "jwt", cfg.Secrets.Jwt, "should be equal")
}

func TestConfigInvalid(t *testing.T) {
	assert.Panics(t, func() {
		config.InitConfig("test_panic")
	}, "should panic")
}
