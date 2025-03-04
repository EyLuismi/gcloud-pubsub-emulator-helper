package pubsub

import (
	"net/http"
	"testing"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubscription(t *testing.T) {
	mockClient := &utils.MockClient{
		ResponseHistory: []utils.MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusNotFound}, Error: nil},
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	err := CreateSubscription(mockClient, "test-project", "projects/test-project/subscriptions/test-subscription", "projects/test-project/topics/test-topic", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mockClient.RequestHistory))
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[0].Method)
	assert.Equal(t, "projects/test-project/subscriptions/test-subscription", mockClient.RequestHistory[0].Path)
	assert.Equal(t, http.MethodPut, mockClient.RequestHistory[1].Method)
	assert.Equal(t, "projects/test-project/subscriptions/test-subscription", mockClient.RequestHistory[1].Path)
}

func TestIsSubscriptionPresent(t *testing.T) {
	mockClient := &utils.MockClient{
		ResponseHistory: []utils.MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	exists, err := IsSubscriptionPresent(mockClient, "test-project", "projects/test-project/subscriptions/test-subscription")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, 1, len(mockClient.RequestHistory))
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[0].Method)
	assert.Equal(t, "projects/test-project/subscriptions/test-subscription", mockClient.RequestHistory[0].Path)
}

func TestListSubscriptions(t *testing.T) {
	mockClient := &utils.MockClient{
		ResponseHistory: []utils.MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK, Body: []byte(`{"subscriptions":[{"name":"test-subscription"}]}`)}, Error: nil},
		},
	}

	subscriptions, err := ListSubscriptions(mockClient, "test-project")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(subscriptions))
	assert.Equal(t, "test-subscription", subscriptions[0].Name)
	assert.Equal(t, 1, len(mockClient.RequestHistory))
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[0].Method)
	assert.Equal(t, "projects/test-project/subscriptions", mockClient.RequestHistory[0].Path)
}

func TestDeleteSubscription(t *testing.T) {
	mockClient := &utils.MockClient{
		ResponseHistory: []utils.MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	err := DeleteSubscription(mockClient, "test-project", "projects/test-project/subscriptions/test-subscription")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockClient.RequestHistory))
	assert.Equal(t, http.MethodDelete, mockClient.RequestHistory[0].Method)
	assert.Equal(t, "projects/test-project/subscriptions/test-subscription", mockClient.RequestHistory[0].Path)
}
