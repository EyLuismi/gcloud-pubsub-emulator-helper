package pubsub

import (
	"net/http"
	"testing"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
	"github.com/stretchr/testify/assert"
)

type MockClientHistoryRequest struct {
	method string
	path   string
	body   []byte
}

type MockClientHistoryResponse struct {
	Response utils.Response
	Error    error
}

type MockClient struct {
	requestHistory  []MockClientHistoryRequest
	responseHistory []MockClientHistoryResponse
}

func (m *MockClient) makeCall(method string, path string, body []byte) (utils.Response, error) {
	m.requestHistory = append(m.requestHistory, MockClientHistoryRequest{
		method: method,
		path:   path,
		body:   body,
	})

	returnValue := m.responseHistory[0]
	m.responseHistory = m.responseHistory[1:]

	return returnValue.Response, returnValue.Error
}

func (m *MockClient) Get(path string) (utils.Response, error) {
	return m.makeCall(http.MethodGet, path, []byte(""))
}

func (m *MockClient) Post(path string, body []byte) (utils.Response, error) {
	return m.makeCall(http.MethodPost, path, body)
}

func (m *MockClient) Put(path string, body []byte) (utils.Response, error) {
	return m.makeCall(http.MethodPut, path, body)
}

func (m *MockClient) Delete(path string) (utils.Response, error) {
	return m.makeCall(http.MethodDelete, path, []byte(""))
}

func (m *MockClient) Patch(path string, body []byte) (utils.Response, error) {
	return m.makeCall(http.MethodPatch, path, body)
}

func TestCreateTopic(t *testing.T) {
	mockClient := &MockClient{
		responseHistory: []MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusNotFound}, Error: nil},
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	err := CreateTopic(mockClient, "test-project", "projects/test-project/topics/test-topic", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mockClient.requestHistory))
	assert.Equal(t, http.MethodGet, mockClient.requestHistory[0].method)
	assert.Equal(t, "projects/test-project/topics/test-topic", mockClient.requestHistory[0].path)
	assert.Equal(t, http.MethodPut, mockClient.requestHistory[1].method)
	assert.Equal(t, "projects/test-project/topics/test-topic", mockClient.requestHistory[1].path)
}

func TestIsTopicPresent(t *testing.T) {
	mockClient := &MockClient{
		responseHistory: []MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	exists, err := IsTopicPresent(mockClient, "projects/test-project/topics/test-topic")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, 1, len(mockClient.requestHistory))
	assert.Equal(t, http.MethodGet, mockClient.requestHistory[0].method)
	assert.Equal(t, "projects/test-project/topics/test-topic", mockClient.requestHistory[0].path)
}

func TestListTopics(t *testing.T) {
	mockClient := &MockClient{
		responseHistory: []MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK, Body: []byte(`{"topics":[{"name":"test-topic"}]}`)}, Error: nil},
		},
	}

	topics, err := ListTopics(mockClient, "test-project")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(topics))
	assert.Equal(t, "test-topic", topics[0].Name)
	assert.Equal(t, 1, len(mockClient.requestHistory))
	assert.Equal(t, http.MethodGet, mockClient.requestHistory[0].method)
	assert.Equal(t, "projects/test-project/topics", mockClient.requestHistory[0].path)
}

func TestDeleteTopic(t *testing.T) {
	mockClient := &MockClient{
		responseHistory: []MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	err := DeleteTopic(mockClient, "test-project", "projects/test-project/topics/test-topic")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockClient.requestHistory))
	assert.Equal(t, http.MethodDelete, mockClient.requestHistory[0].method)
	assert.Equal(t, "projects/test-project/topics/test-topic", mockClient.requestHistory[0].path)
}
