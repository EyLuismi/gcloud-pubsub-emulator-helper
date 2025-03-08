package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
)

// https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics#MessageStoragePolicy
type TopicMessageStoragePolicy struct {
	AllowedPersistenceRegions []string `json:"allowedPersistenceRegions"`
	EnforceInTransit          bool     `json:"enforceInTransit"`
}

// https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics#TopicIngestionDataSourceSettings
type TopicIngestionDataSourceSettings struct {
	PlatformLogsSettings TopicIngestionDataSourceSettingsPlatformLogsSettings `json:"platformLogsSettings"`
	AwsKinesis           *TopicIngestionDataSourceSettingsAwsKinesis          `json:"awsKinesis,omitempty"`
	CloudStorage         *TopicIngestionDataSourceSettingsCloudStorage        `json:"cloudStorage,omitempty"`
}

// https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics#CloudStorage
type TopicIngestionDataSourceSettingsCloudStorage struct {
	State                   CloudStorageState       `json:"state"`
	Bucket                  string                  `json:"bucket"`
	MinimumObjectCreateTime string                  `json:"minimumObjectCreateTime"`
	MatchGlob               string                  `json:"matchGlob"`
	TextFormat              *CloudStorageTextFormat `json:"textFormat,omitempty"`
	AvroFormat              *struct{}               `json:"avroFormat,omitempty"`
	PubSubAvroFormat        *struct{}               `json:"pubsubAvroFormat,omitempty"`
}

type CloudStorageState string

const (
	CLOUD_STORAGE_STATE_UNSPECIFIED               CloudStorageState = "STATE_UNSPECIFIED"
	CLOUD_STORAGE_ACTIVE                          CloudStorageState = "ACTIVE"
	CLOUD_STORAGE_CLOUD_STORAGE_PERMISSION_DENIED CloudStorageState = "CLOUD_STORAGE_PERMISSION_DENIED"
	CLOUD_STORAGE_PUBLISH_PERMISSION_DENIED       CloudStorageState = "PUBLISH_PERMISSION_DENIED"
	CLOUD_STORAGE_BUCKET_NOT_FOUND                CloudStorageState = "BUCKET_NOT_FOUND"
	CLOUD_STORAGE_TOO_MANY_OBJECTS                CloudStorageState = "TOO_MANY_OBJECTS"
)

type CloudStorageTextFormat struct {
	Delimiter string `json:"delimiter"`
}

// https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics#PlatformLogsSettings
type TopicIngestionDataSourceSettingsPlatformLogsSettings struct {
	Severity string `json:"severity"`
}

type AwsKinesisState string

const (
	AWS_KINESIS_STATE_UNSPECIFIED         AwsKinesisState = "STATE_UNSPECIFIED"
	AWS_KINESIS_ACTIVE                    AwsKinesisState = "ACTIVE"
	AWS_KINESIS_KINESIS_PERMISSION_DENIED AwsKinesisState = "KINESIS_PERMISSION_DENIED"
	AWS_KINESIS_PUBLISH_PERMISSION_DENIED AwsKinesisState = "PUBLISH_PERMISSION_DENIED"
	AWS_KINESIS_STREAM_NOT_FOUND          AwsKinesisState = "STREAM_NOT_FOUND"
	AWS_KINESIS_CONSUMER_NOT_FOUND        AwsKinesisState = "CONSUMER_NOT_FOUND"
)

// https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics#AwsKinesis
type TopicIngestionDataSourceSettingsAwsKinesis struct {
	State             AwsKinesisState `json:"state"`
	StreamArn         string          `json:"streamArn"`
	ConsumerArn       string          `json:"consumerArn"`
	AwsRoleArn        string          `json:"awsRoleArn"`
	GcpServiceAccount string          `json:"gcpServiceAccount"`
}

// Topic represents a Pub/Sub topic.
type Topic struct {
	Name                 string                    `json:"name"`
	Subscriptions        []Subscription            `json:"subscriptions"`
	Labels               Labels                    `json:"labels"`
	MessageStoragePolicy TopicMessageStoragePolicy `json:"messageStoragePolicy"`

	// https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics#state
	/**
	  Output only. An output-only field indicating the state of the topic.
	  => The emulator seems to print an empty value
	*/
	State string `json:"state"`

	KmsKeyName string `json:"kmsKeyName"`

	/**
	  This field is not accepted by the emulator, but have been added as the
	    official REST API accepts it. You shouldn't send it in the mean time.
	*/
	MessageRetentionDuration string `json:"messageRetentionDuration"`

	IngestionDataSourceSettings *TopicIngestionDataSourceSettings `json:"ingestionDataSourceSettings,omitempty"`
}

// String returns a JSON string representation of the Topic.
func (t *Topic) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("error marshaling topic: %v", err)
	}
	return string(b)
}

// CreateTopic creates a topic if it does not exist.
func CreateTopic(
	client utils.ClientInterface,
	project, topicResourceName string,
	labels *Labels,
	messageStoragePolicy *TopicMessageStoragePolicy,
	kmsKeyName string,
	messageRetentionDuration string,
	ingestionDataSourceSettings *TopicIngestionDataSourceSettings,
) error {
	// Check if the topic already exists.
	exists, err := IsTopicPresent(client, topicResourceName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	type CreateTopicBody struct {
		Labels                      Labels                            `json:"labels"`
		MessageStoragePolicy        TopicMessageStoragePolicy         `json:"messageStoragePolicy"`
		KmsKeyName                  string                            `json:"kmsKeyName"`
		MessageRetentionDuration    *string                           `json:"messageRetentionDuration"`
		IngestionDataSourceSettings *TopicIngestionDataSourceSettings `json:"ingestionDataSourceSettings"`
	}

	createTopicBody := CreateTopicBody{}

	if labels != nil {
		createTopicBody.Labels = *labels
	}

	if messageStoragePolicy != nil {
		createTopicBody.MessageStoragePolicy = *messageStoragePolicy
	}

	if kmsKeyName != "" {
		createTopicBody.KmsKeyName = kmsKeyName
	}

	if messageRetentionDuration != "" {
		createTopicBody.MessageRetentionDuration = &messageRetentionDuration
	} else {
		createTopicBody.MessageRetentionDuration = nil
	}

	if ingestionDataSourceSettings != nil {
		createTopicBody.IngestionDataSourceSettings = ingestionDataSourceSettings
	}

	jsonCreateTopicBody, err := json.Marshal(createTopicBody)
	if err != nil {
		return err
	}

	// Create the topic with an empty configuration.
	response, err := client.Put(topicResourceName, jsonCreateTopicBody)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating topic: status code %d", response.StatusCode)
	}

	return nil
}

// IsTopicPresent checks if a topic exists.
func IsTopicPresent(
	client utils.ClientInterface,
	topicResourceName string,
) (bool, error) {
	response, err := client.Get(topicResourceName)
	if err != nil {
		return false, err
	}

	switch response.StatusCode {
	case http.StatusNotFound:
		return false, nil
	case http.StatusOK:
		return true, nil
	default:
		return false, fmt.Errorf("unexpected status code %d in IsTopicPresent", response.StatusCode)
	}
}

// ListTopics lists all topics of a project.
func ListTopics(
	client utils.ClientInterface,
	project string,
) ([]Topic, error) {
	// Build the URL for listing topics.
	url := fmt.Sprintf("projects/%s/topics", project)
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusNotFound:
		return nil, errors.New("project not found")
	case http.StatusOK:
		var res listTopicsResponse
		if err := json.Unmarshal(response.Body, &res); err != nil {
			return nil, err
		}
		return res.Topics, nil
	default:
		return nil, fmt.Errorf("unexpected status code %d in ListTopics", response.StatusCode)
	}
}

// listTopicsResponse is the internal structure for unmarshalling the ListTopics response.
type listTopicsResponse struct {
	Topics []Topic `json:"topics"`
}

// DeleteTopic deletes a topic.
func DeleteTopic(
	client utils.ClientInterface,
	project, topicResourceName string,
) error {
	response, err := client.Delete(topicResourceName)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting topic: status code %d", response.StatusCode)
	}

	return nil
}

// GetResourceNameForTopic generates the full resource name for a topic.
func GetResourceNameForTopic(project, topic string) string {
	return fmt.Sprintf("projects/%s/topics/%s", project, topic)
}
