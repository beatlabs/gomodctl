package godoc

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Query(t *testing.T) {
	client := NewClient(context.TODO())

	response, err := client.Search("mock")

	assert.NoError(t, err)
	assert.NotNil(t, response)

	for _, result := range response {
		log.Printf("%v\n", result)
	}
}

func TestClient_Info(t *testing.T) {
	client := NewClient(context.TODO())

	response, err := client.Info("patates")

	assert.Error(t, err)
	assert.Empty(t, response)
}

func TestClient_Info2(t *testing.T) {
	client := NewClient(context.TODO())

	response, err := client.Info("github.com/stretchr/testify/mock")

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}
