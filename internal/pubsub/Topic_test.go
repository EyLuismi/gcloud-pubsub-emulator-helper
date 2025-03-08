package pubsub

import (
	"net/http"
	"testing"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_Topics_Create(t *testing.T) {
	mockClient := &utils.MockClient{
		ResponseHistory: []utils.MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusNotFound}, Error: nil},
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	err := CreateTopic(mockClient, "test-project", "projects/test-project/topics/test-topic", nil, nil, "")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mockClient.RequestHistory))
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[0].Method)
	assert.Equal(t, "projects/test-project/topics/test-topic", mockClient.RequestHistory[0].Path)
	assert.Equal(t, http.MethodPut, mockClient.RequestHistory[1].Method)
	assert.Equal(t, "projects/test-project/topics/test-topic", mockClient.RequestHistory[1].Path)
}

func Test_Topics_IsPresent(t *testing.T) {
	mockClient := &utils.MockClient{
		ResponseHistory: []utils.MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	exists, err := IsTopicPresent(mockClient, "projects/test-project/topics/test-topic")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, 1, len(mockClient.RequestHistory))
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[0].Method)
	assert.Equal(t, "projects/test-project/topics/test-topic", mockClient.RequestHistory[0].Path)
}

func Test_Topics_List(t *testing.T) {
	mockClient := &utils.MockClient{
		ResponseHistory: []utils.MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK, Body: []byte(`{"topics":[{"name":"test-topic"}]}`)}, Error: nil},
		},
	}

	topics, err := ListTopics(mockClient, "test-project")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(topics))
	assert.Equal(t, "test-topic", topics[0].Name)
	assert.Equal(t, 1, len(mockClient.RequestHistory))
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[0].Method)
	assert.Equal(t, "projects/test-project/topics", mockClient.RequestHistory[0].Path)
}

func Test_Topics_Delete(t *testing.T) {
	mockClient := &utils.MockClient{
		ResponseHistory: []utils.MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	err := DeleteTopic(mockClient, "test-project", "projects/test-project/topics/test-topic")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockClient.RequestHistory))
	assert.Equal(t, http.MethodDelete, mockClient.RequestHistory[0].Method)
	assert.Equal(t, "projects/test-project/topics/test-topic", mockClient.RequestHistory[0].Path)
}
