package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	rawContent := `
rpc:
- name: hello
  endpoint: https://hello.world
- name: world
  endpoint: https://dsrv.sui
`
	c, err := parseConfig([]byte(rawContent))

	assert.Nil(t, err)
	assert.Equal(t, c.RPC[0].Name, "hello")
	assert.Equal(t, c.RPC[0].Endpoint, "https://hello.world")

	assert.Equal(t, c.RPC[1].Name, "world")
	assert.Equal(t, c.RPC[1].Endpoint, "https://dsrv.sui")
}

func TestConfigLoad(t *testing.T) {
	c, err := Load("./fixture/test_config.yaml")

	assert.Nil(t, err)
	assert.Equal(t, c.RPC[0].Name, "hello")
	assert.Equal(t, c.RPC[0].Endpoint, "https://hello.world")

	assert.Equal(t, c.RPC[1].Name, "world")
	assert.Equal(t, c.RPC[1].Endpoint, "https://dsrv.sui")
}
