package internal

import (
	"net/http"
	"testing"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/pubsub"
	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_Configuration_LoadFile_WithoutProjects(t *testing.T) {
	filepath := "test_config.json"
	mockReader := utils.NewFileReaderMockBasic(
		`{
      "host": "localhost:8888",
      "startTimeoutMs": 30000,
      "avoidStartupCheck": false,
      "projects": [],
      "timeBetweenStartupChecksMs": 200,
      "delayBeforeStartupCheckMs": 0
    }`,
	)

	config, err := LoadConfigurationFromFile(mockReader, filepath)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8888", config.Host)
	assert.Equal(t, 30000, config.StartTimeoutMs)
	assert.False(t, config.AvoidStartupCheck)
	assert.Equal(t, 200, config.TimeBetweenStartupChecksMs)
	assert.Equal(t, 0, config.DelayBeforeStartupCheckMs)
	assert.Equal(t, 0, len(config.Projects))
}

func Test_Configuration_LoadFile_WithProjectWithoutTopics(t *testing.T) {
	filepath := "test_config.json"
	mockReader := utils.NewFileReaderMockBasic(
		`{
      "host": "localhost:8888",
      "startTimeoutMs": 30000,
      "avoidStartupCheck": false,
      "projects": [{
        "name": "first-project"
      }],
      "timeBetweenStartupChecksMs": 200,
      "delayBeforeStartupCheckMs": 0
    }`,
	)

	config, err := LoadConfigurationFromFile(mockReader, filepath)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8888", config.Host)
	assert.Equal(t, 30000, config.StartTimeoutMs)
	assert.False(t, config.AvoidStartupCheck)
	assert.Equal(t, 200, config.TimeBetweenStartupChecksMs)
	assert.Equal(t, 0, config.DelayBeforeStartupCheckMs)
	assert.Equal(t, 1, len(config.Projects))
	assert.Equal(t, 0, len(config.Projects[0].Topics))
}

func Test_Configuration_LoadFile_WithProjectWithTopicWithoutSubscriptions(t *testing.T) {
	filepath := "test_config.json"
	mockReader := utils.NewFileReaderMockBasic(
		`{
      "host": "localhost:8888",
      "startTimeoutMs": 30000,
      "avoidStartupCheck": false,
      "projects": [{
        "name": "first-project",
        "topics": [
          {
            "name": "testing.new-topic.v1"
          }
        ]
      }],
      "timeBetweenStartupChecksMs": 200,
      "delayBeforeStartupCheckMs": 0
    }`,
	)

	config, err := LoadConfigurationFromFile(mockReader, filepath)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8888", config.Host)
	assert.Equal(t, 30000, config.StartTimeoutMs)
	assert.False(t, config.AvoidStartupCheck)
	assert.Equal(t, 200, config.TimeBetweenStartupChecksMs)
	assert.Equal(t, 0, config.DelayBeforeStartupCheckMs)
	assert.Equal(t, 1, len(config.Projects))
	assert.Equal(t, 1, len(config.Projects[0].Topics))
	assert.Equal(t, 0, len(config.Projects[0].Topics[0].Subscriptions))
}

func Test_Configuration_LoadFile_WithProjectWithTopicWithSubscription(t *testing.T) {
	filepath := "test_config.json"
	mockReader := utils.NewFileReaderMockBasic(
		`{
      "host": "localhost:8888",
      "startTimeoutMs": 30000,
      "avoidStartupCheck": false,
      "projects": [{
        "name": "first-project",
        "topics": [
          {
            "name": "testing.new-topic.v1",
            "subscriptions": [
              {
                "name": "testing.new-topic.v1.subscripion1"
              }
            ]
          }
        ]
      }],
      "timeBetweenStartupChecksMs": 200,
      "delayBeforeStartupCheckMs": 0
    }`,
	)

	config, err := LoadConfigurationFromFile(mockReader, filepath)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8888", config.Host)
	assert.Equal(t, 30000, config.StartTimeoutMs)
	assert.False(t, config.AvoidStartupCheck)
	assert.Equal(t, 200, config.TimeBetweenStartupChecksMs)
	assert.Equal(t, 0, config.DelayBeforeStartupCheckMs)
	assert.Equal(t, 1, len(config.Projects))
	assert.Equal(t, 1, len(config.Projects[0].Topics))
	assert.Equal(t, 1, len(config.Projects[0].Topics[0].Subscriptions))
}

func Test_Configuration_LoadFile_WithProjectWithTopicWithKMSKeyNameWithSubscription(t *testing.T) {
	filepath := "test_config.json"
	mockReader := utils.NewFileReaderMockBasic(
		`{
      "host": "localhost:8888",
      "startTimeoutMs": 30000,
      "avoidStartupCheck": false,
      "projects": [{
        "name": "first-project",
        "topics": [
          {
            "name": "testing.new-topic.v1",
            "kmsKeyName": "projects/advanced-configuration-example/locations/europe-west2/keyRings/advanced-configuration-example/cryptoKeys/advanced-configuration-example-key",
            "subscriptions": [
              {
                "name": "testing.new-topic.v1.subscripion1"
              }
            ]
          }
        ]
      }],
      "timeBetweenStartupChecksMs": 200,
      "delayBeforeStartupCheckMs": 0
    }`,
	)

	config, err := LoadConfigurationFromFile(mockReader, filepath)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8888", config.Host)
	assert.Equal(t, 30000, config.StartTimeoutMs)
	assert.False(t, config.AvoidStartupCheck)
	assert.Equal(t, 200, config.TimeBetweenStartupChecksMs)
	assert.Equal(t, 0, config.DelayBeforeStartupCheckMs)
	assert.Equal(t, 1, len(config.Projects))
	assert.Equal(t, 1, len(config.Projects[0].Topics))
	assert.Equal(t, 1, len(config.Projects[0].Topics[0].Subscriptions))
	assert.Equal(t, "projects/advanced-configuration-example/locations/europe-west2/keyRings/advanced-configuration-example/cryptoKeys/advanced-configuration-example-key", config.Projects[0].Topics[0].KmsKeyName)
}

func Test_Configuration_LoadFile_WithProjectWithTopicWithEmptyLabelsAndWithSubscription(t *testing.T) {
	filepath := "test_config.json"
	mockReader := utils.NewFileReaderMockBasic(
		`{
      "host": "localhost:8888",
      "startTimeoutMs": 30000,
      "avoidStartupCheck": false,
      "projects": [{
        "name": "first-project",
        "topics": [
          {
            "name": "testing.new-topic.v1",
            "labels": {},
            "subscriptions": [
              {
                "name": "testing.new-topic.v1.subscripion1"
              }
            ]
          }
        ]
      }],
      "timeBetweenStartupChecksMs": 200,
      "delayBeforeStartupCheckMs": 0
    }`,
	)

	config, err := LoadConfigurationFromFile(mockReader, filepath)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8888", config.Host)
	assert.Equal(t, 30000, config.StartTimeoutMs)
	assert.False(t, config.AvoidStartupCheck)
	assert.Equal(t, 200, config.TimeBetweenStartupChecksMs)
	assert.Equal(t, 0, config.DelayBeforeStartupCheckMs)
	assert.Equal(t, 1, len(config.Projects))
	assert.Equal(t, 1, len(config.Projects[0].Topics))
	assert.Equal(t, 1, len(config.Projects[0].Topics[0].Subscriptions))
	assert.Equal(t, 0, len(config.Projects[0].Topics[0].Labels))
}

func Test_Configuration_LoadFile_WithProjectWithTopicWithALabelAndWithSubscription(t *testing.T) {
	filepath := "test_config.json"
	mockReader := utils.NewFileReaderMockBasic(
		`{
      "host": "localhost:8888",
      "startTimeoutMs": 30000,
      "avoidStartupCheck": false,
      "projects": [{
        "name": "first-project",
        "topics": [
          {
            "name": "testing.new-topic.v1",
            "labels": {
              "firstLabel": "labelValue"
            },
            "subscriptions": [
              {
                "name": "testing.new-topic.v1.subscripion1"
              }
            ]
          }
        ]
      }],
      "timeBetweenStartupChecksMs": 200,
      "delayBeforeStartupCheckMs": 0
    }`,
	)

	config, err := LoadConfigurationFromFile(mockReader, filepath)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8888", config.Host)
	assert.Equal(t, 30000, config.StartTimeoutMs)
	assert.False(t, config.AvoidStartupCheck)
	assert.Equal(t, 200, config.TimeBetweenStartupChecksMs)
	assert.Equal(t, 0, config.DelayBeforeStartupCheckMs)
	assert.Equal(t, 1, len(config.Projects))
	assert.Equal(t, 1, len(config.Projects[0].Topics))
	assert.Equal(t, 1, len(config.Projects[0].Topics[0].Subscriptions))
	assert.Equal(t, 1, len(config.Projects[0].Topics[0].Labels))
}

func Test_Configuration_LoadFile_WithProjectWithTopicWithTwoLabelsTheSameAndWithSubscription(t *testing.T) {
	filepath := "test_config.json"
	mockReader := utils.NewFileReaderMockBasic(
		`{
      "host": "localhost:8888",
      "startTimeoutMs": 30000,
      "avoidStartupCheck": false,
      "projects": [{
        "name": "first-project",
        "topics": [
          {
            "name": "testing.new-topic.v1",
            "labels": {
              "firstLabel": "labelValue1",
              "firstLabel": "labelValue2"
            },
            "subscriptions": [
              {
                "name": "testing.new-topic.v1.subscripion1"
              }
            ]
          }
        ]
      }],
      "timeBetweenStartupChecksMs": 200,
      "delayBeforeStartupCheckMs": 0
    }`,
	)

	config, err := LoadConfigurationFromFile(mockReader, filepath)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8888", config.Host)
	assert.Equal(t, 30000, config.StartTimeoutMs)
	assert.False(t, config.AvoidStartupCheck)
	assert.Equal(t, 200, config.TimeBetweenStartupChecksMs)
	assert.Equal(t, 0, config.DelayBeforeStartupCheckMs)
	assert.Equal(t, 1, len(config.Projects))
	assert.Equal(t, 1, len(config.Projects[0].Topics))
	assert.Equal(t, 1, len(config.Projects[0].Topics[0].Subscriptions))
	assert.Equal(t, 1, len(config.Projects[0].Topics[0].Labels))
	assert.Equal(t, "labelValue2", config.Projects[0].Topics[0].Labels["firstLabel"])
}

func Test_Configuration_ReplaceHost(t *testing.T) {
	config := Configuration{Host: "localhost:8085"}
	newHost := "0.0.0.0:8085"
	config = config.ReplaceHost(newHost)
	assert.Equal(t, newHost, config.Host)
}

func Test_Configuration_Sync(t *testing.T) {
	mockClient := &utils.MockClient{
		ResponseHistory: []utils.MockClientHistoryResponse{
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
			{Response: utils.Response{StatusCode: http.StatusOK, Body: []byte(`{"topics":[]}`)}, Error: nil},
			{Response: utils.Response{StatusCode: http.StatusOK, Body: []byte(`{"subscriptions":[]}`)}, Error: nil},
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
			{Response: utils.Response{StatusCode: http.StatusOK}, Error: nil},
		},
	}

	config := Configuration{
		Host: "localhost:8085",
		Projects: []pubsub.Project{
			{
				Name: "test-project",
				Topics: []pubsub.Topic{
					{
						Name: "test-topic",
						Subscriptions: []pubsub.Subscription{
							{Name: "test-subscription"},
						},
					},
				},
			},
		},
	}

	config.Sync(mockClient)
	assert.Equal(t, 5, len(mockClient.RequestHistory))
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[0].Method)
	assert.Equal(t, "", mockClient.RequestHistory[0].Path)
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[1].Method)
	assert.Equal(t, "projects/test-project/topics", mockClient.RequestHistory[1].Path)
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[2].Method)
	assert.Equal(t, "projects/test-project/subscriptions", mockClient.RequestHistory[2].Path)
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[3].Method)
	assert.Equal(t, "projects/test-project/topics/test-topic", mockClient.RequestHistory[3].Path)
	assert.Equal(t, http.MethodGet, mockClient.RequestHistory[4].Method)
	assert.Equal(t, "projects/test-project/subscriptions/test-subscription", mockClient.RequestHistory[4].Path)
}
