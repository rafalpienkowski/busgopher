package controller

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

func getTestConfig() *config.InMemoryConfigStorage {

	return &config.InMemoryConfigStorage{
		Config: config.Config{
			Connections: []asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: []asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		},
	}
}

func TestControllerShouldSetLoadedConfig(t *testing.T) {
    testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig

	sut, _ := NewController(givenConfig)

	assert.Equal(t, testConfig.Config, sut.Config)
}

func TestControllerShouldSetConnectionByName(t *testing.T) {
    testConfig := getTestConfig()
	var givenConfig config.ConfigStorage = testConfig
    var buffer bytes.Buffer
	var writer io.Writer = &buffer

	sut, _ := NewController(givenConfig)
    sut.SetLogsWriter(writer)
    sut.SelectConnectionByName("test")

	assert.Equal(t, &(testConfig.Config.Connections)[0], sut.SelectedConnection)
}
