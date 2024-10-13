package controller

import ( "testing"
    "github.com/stretchr/testify/assert"

	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
)

func TestControllerShouldSetLoadedConnection(t *testing.T) {
    givenConfig := &config.InMemoryConfigStorage{
		Config: config.Config{
			Connections: &[]asb.Connection{
				{
					Name:      "test",
					Namespace: "test.azure.com",
					Destinations: []string{
						"queue",
						"topic",
					},
				},
			},
			Messages: &[]asb.Message{
				{
					Name: "test",
					Body: "test msg body",
				},
			},
		},
	}
    var gc config.ConfigStorage = givenConfig

    sut,err := NewController(gc)
    if err != nil {
        t.Errorf("Can't create controller: %v", err.Error())
    }

    assert.Equal(t, givenConfig.Config.Connections, &sut.Connections)
    assert.Equal(t, givenConfig.Config.Messages, &sut.Messages)
}
