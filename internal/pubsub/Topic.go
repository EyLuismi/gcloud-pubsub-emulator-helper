package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/EyLuismi/gcloud-pubsub-emulator-helper/internal/utils"
)

// Topic represents a Pub/Sub topic.
type Topic struct {
	Name          string         `json:"name"`
	Subscriptions []Subscription `json:"subscriptions"`
	Labels        Labels         `json:"labels"`
}

// CreateTopic creates a topic if it does not exist.
func CreateTopic(
	client utils.Client,
	project, topicResourceName string,
	labels *Labels,
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
		Labels Labels `json:"labels"`
	}

	createTopicBody := CreateTopicBody{}

	if labels != nil {
		createTopicBody.Labels = *labels
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
	client utils.Client,
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
	client utils.Client,
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
	client utils.Client,
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
